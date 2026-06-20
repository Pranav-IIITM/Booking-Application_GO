package config

import (
	"context"
	"errors"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Config struct {
	Port                    string
	FirebaseCredentialsPath string
	FirebaseAuth            *auth.Client
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:                    getEnv("PORT", "8080"),
		FirebaseCredentialsPath: os.Getenv("FIREBASE_CREDENTIALS_PATH"),
	}

	if cfg.FirebaseCredentialsPath == "" {
		return nil, errors.New("FIREBASE_CREDENTIALS_PATH is required")
	}

	return cfg, nil
}

func InitFirebase() (*auth.Client, *firestore.Client, error) {
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credentialsPath == "" {
		return nil, nil, errors.New("FIREBASE_CREDENTIALS_PATH is required")
	}

	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		return nil, nil, err
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, nil, err
	}

	firestoreClient, err := app.Firestore(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return authClient, firestoreClient, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
