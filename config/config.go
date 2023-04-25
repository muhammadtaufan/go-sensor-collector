package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GRPC_HOST         string
	GRPC_PORT         string
	DATABASE_HOST     string
	DATABASE_USERNAME string
	DATABASE_PASSWORD string
	DATABASE_NAME     string
	API_HOST          string
	API_PORT          string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	return &Config{
		GRPC_HOST:         getEnv("GRPC_HOST", "0.0.0.0"),
		GRPC_PORT:         getEnv("GRPC_PORT", "50051"),
		API_HOST:          getEnv("API_HOST", "0.0.0.0"),
		API_PORT:          getEnv("API_PORT", "3001"),
		DATABASE_HOST:     getEnv("DATABASE_HOST", "localhost"),
		DATABASE_USERNAME: getEnv("DATABASE_USERNAME", "mysql"),
		DATABASE_PASSWORD: getEnv("DATABASE_PASSWORD", "mysql"),
		DATABASE_NAME:     getEnv("DATABASE_NAME", "sensor_collector"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
