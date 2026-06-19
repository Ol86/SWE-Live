package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Port:        getEnv("PORT", ":8080"),
		Environment: getEnv("ENVIRONMENT", "production"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://library:p@localhost:5432/library?sslmode=disable"),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		if key == "PORT" {
			return ":" + value
		}
		return value
	}
	return fallback
}
