package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

var config *Config

func loadConfig() *Config {
	// Load .env file if it exists.
	// Do NOT use log.Fatal here because in Docker/Production, environment variables
	// are injected directly by docker-compose without needing a physical .env file inside the container.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables from system/Docker")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "pendataan_apar_db"),
	}
}

func GetConfig() *Config {
	if config == nil {
		config = loadConfig()
	}

	return config
}

// Helper function to read an env variable or return a fallback default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}