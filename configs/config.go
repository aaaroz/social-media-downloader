package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Environment   string
	MaxDownloadMB int
}

func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := getEnv("PORT", "3000")
	env := getEnv("ENVIRONMENT", "development")
	maxDownloadMB, _ := strconv.Atoi(getEnv("MAX_DOWNLOAD_MB", "100"))

	return &Config{
		Port:          port,
		Environment:   env,
		MaxDownloadMB: maxDownloadMB,
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
