package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host        string
	DatabaseURL string
	JWTSecret   string
}

func Load() (Config, error) {
	godotenv.Load()

	host := os.Getenv("APP_HOST")

	if host == "" {
		return Config{}, NewEnvError("APP_HOST")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		return Config{}, NewEnvError("DATABASE_URL")
	}

	JWTSecret := os.Getenv("JWT_SECRET")

	if JWTSecret == "" {
		return Config{}, NewEnvError("JWT_SECRET")
	}

	return Config{
		Host:        host,
		DatabaseURL: databaseURL,
		JWTSecret:   JWTSecret,
	}, nil
}
