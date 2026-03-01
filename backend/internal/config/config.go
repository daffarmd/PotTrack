package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	DBUrl     string
	JWTSecret string
	Port      string
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	db := os.Getenv("DATABASE_URL")
	if db == "" {
		db = "postgres://postgres:password@localhost:5432/pot_track_dev?sslmode=disable"
	}
	jwt := os.Getenv("JWT_SECRET")
	if jwt == "" {
		jwt = "secret"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		DBUrl:     db,
		JWTSecret: jwt,
		Port:      port,
	}, nil
}
