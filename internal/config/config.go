package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	MongoName string
	JWTSecret string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("Failed to load .env: %w", err)
	}

	cfg := Config{
		MongoURI:  strings.TrimSpace(os.Getenv("MONGODB_URI")),
		MongoName: strings.TrimSpace(os.Getenv("MONGODB_NAME")),
		JWTSecret: strings.TrimSpace(os.Getenv("JWT_SECRET")),
	}

	if cfg.MongoURI == "" {
		return Config{}, fmt.Errorf("Missing Mongo URI")
	}
	if cfg.MongoName == "" {
		return Config{}, fmt.Errorf("Missing Mongo Name")
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("Missing Mongo JWT Secret")
	}

	return cfg, nil

}
