package main

import (
	"ad-tracking-system/internal/api"
	"ad-tracking-system/internal/config"
	"ad-tracking-system/internal/domain/services"
	"ad-tracking-system/internal/repository"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv" // Add this import
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Enable debug mode if DEBUG environment variable is set
	debugMode := os.Getenv("DEBUG") == "true"
	if debugMode {
		log.Println("Debug mode enabled")
	}

	// Log configuration for debugging
	if debugMode {
		log.Printf("Loaded configuration: %+v", cfg)
	}

	// Initialize the database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	log.Println("Successfully connected to the database")

	// Initialize repositories
	adRepo := repository.NewAdRepository(db)

	// Initialize services
	adService := services.NewAdService(adRepo)

	// Initialize the API router
	router := api.NewRouter(adService)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.HTTPPort), // Convert HTTPPort to string
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Start the HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %d", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}
	log.Println("HTTP server stopped")

	// Close database connection
	if err := db.Close(); err != nil {
		log.Fatalf("Database shutdown error: %v", err)
	}
	log.Println("Database connection closed")

	log.Println("Server shutdown complete")
}
