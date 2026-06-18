package config

import (
	"context"
	"errors"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Config struct {
	Port                    string
	DatabaseURL             string
	FirebaseCredentialsPath string
	FirebaseAuth            *auth.Client
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:                    getEnv("PORT", "8080"),
		DatabaseURL:             os.Getenv("DATABASE_URL"),
		FirebaseCredentialsPath: os.Getenv("FIREBASE_CREDENTIALS_PATH"),
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if cfg.FirebaseCredentialsPath == "" {
		return nil, errors.New("FIREBASE_CREDENTIALS_PATH is required")
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(cfg.FirebaseCredentialsPath))
	if err != nil {
		return nil, err
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	cfg.FirebaseAuth = authClient

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
