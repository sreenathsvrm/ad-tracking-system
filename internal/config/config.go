package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds the application configuration
type Config struct {
	HTTPPort     int
	KafkaBrokers []string
	KafkaTopic   string
	RedisURL     string
	DatabaseURL  string
	MetricsPort  int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Constants for default values
const (
	defaultHTTPPort     = 8080
	defaultMetricsPort  = 2112
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultKafkaTopic   = "ad-clicks"
	defaultRedisURL     = "localhost:6379"
	defaultDatabaseURL  = "postgres://postgres:password@localhost:5432/advertisements?sslmode=disable"
	defaultKafkaBrokers = "localhost:9092"
)

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{
		HTTPPort:     getEnvAsInt("HTTP_PORT", defaultHTTPPort),
		KafkaBrokers: getEnvAsSlice("KAFKA_BROKERS", []string{defaultKafkaBrokers}, ","),
		KafkaTopic:   getEnv("KAFKA_TOPIC", defaultKafkaTopic),
		RedisURL:     getEnv("REDIS_URL", defaultRedisURL),
		DatabaseURL:  getEnv("DATABASE_URL", defaultDatabaseURL),
		MetricsPort:  getEnvAsInt("METRICS_PORT", defaultMetricsPort),
		ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", defaultReadTimeout),   
		WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", defaultWriteTimeout), 
	}

	// Validate critical configurations
	if !strings.HasPrefix(cfg.DatabaseURL, "postgres://") {
		log.Fatal("Invalid DATABASE_URL format")
	}

	// Log the loaded configuration
	log.Printf("Loaded configuration: %+v", cfg)

	return cfg
}

// Helper functions to read environment variables
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getEnvAsSlice(key string, defaultValue []string, separator string) []string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return strings.Split(value, separator)
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}