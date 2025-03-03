package main

import (
	"ad-tracking-system/internal/api"
	"ad-tracking-system/internal/config"
	"ad-tracking-system/internal/domain/services"
	"ad-tracking-system/internal/repository"
	"ad-tracking-system/pkg/kafka"
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Enable debug mode if DEBUG environment variable is set
	debugMode := os.Getenv("DEBUG") == "true"

	// Initialize structured logging
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	if debugMode {
		logger.Info("Debug mode enabled")
	}

	// Log configuration for debugging
	if debugMode {
		logger.Info("Loaded configuration", "config", cfg)
	}

	// Initialize the database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping the database", "error", err)
		os.Exit(1)
	}
	logger.Info("Successfully connected to the database")

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// Test Redis connection
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	logger.Info("Successfully connected to Redis")

	// Initialize repositories
	adRepo := repository.NewAdRepository(db)
	clickRepo := repository.NewClickRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(redisClient)

	// Initialize services
	adService := services.NewAdService(adRepo)
	clickService := services.NewClickService(clickRepo, analyticsRepo)

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		logger.Error("Failed to create Kafka producer", "error", err)
		os.Exit(1)
	}
	defer kafkaProducer.Close()
	logger.Info("Kafka producer initialized")

	// Initialize the API router
	router := api.NewRouter(adService, clickService)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.HTTPPort), // Convert HTTPPort to string
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Start the HTTP server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", "port", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the HTTP server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("HTTP server shutdown error", "error", err)
	}
	logger.Info("HTTP server stopped")

	// Close Kafka producer
	if err := kafkaProducer.Close(); err != nil {
		logger.Error("Kafka producer shutdown error", "error", err)
	}
	logger.Info("Kafka producer stopped")

	// Close Redis connection
	if err := redisClient.Close(); err != nil {
		logger.Error("Redis shutdown error", "error", err)
	}
	logger.Info("Redis connection closed")

	// Close database connection
	if err := db.Close(); err != nil {
		logger.Error("Database shutdown error", "error", err)
	}
	logger.Info("Database connection closed")

	logger.Info("Server shutdown complete")
}