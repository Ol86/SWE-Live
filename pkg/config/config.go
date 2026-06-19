package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Environment string
	// Database connection
	DatabaseURL string
	// TLS
	TLSEnabled  bool
	TLSCertPath string
	TLSKeyPath  string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Port:        getEnv("PORT", ":8080"),
		Environment: getEnv("ENVIRONMENT", "production"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://library:p@localhost:5432/library?sslmode=disable"),
		TLSEnabled:  getBoolEnv("TLS_ENABLED", false),
		TLSCertPath: getEnv("TLS_CERT_PATH", "pkg/config/tls/certificate.crt"),
		TLSKeyPath:  getEnv("TLS_KEY_PATH", "pkg/config/tls/key.pem"),
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

func getBoolEnv(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true"
	}
	return fallback
}
