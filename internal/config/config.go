package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	ServerAddr  string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://admin:root@localhost:5432/scholarly?sslmode=disable"),
		ServerAddr:  getEnv("SERVER_ADDR", "localhost:8080"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
