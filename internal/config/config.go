package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL string
	JWTSecret   string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("Failed to load .env: %w", err)
	}

	cfg := Config{
		PostgresURL: strings.TrimSpace(os.Getenv("POSTGRES_URL")),
		JWTSecret:   strings.TrimSpace(os.Getenv("JWT_SECRET")),
	}

	if cfg.PostgresURL == "" {
		return Config{}, fmt.Errorf("Missing Postgres URL (POSTGRES_URL)")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("Missing JWT secret (JWT_SECRET)")
	}

	return cfg, nil

}
