package config

import (
	"errors"
	"os"
)

type AppConfig struct {
	DatabaseUrl string
	Port        string
	JWTSecret   string
}

func Load() (*AppConfig, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	port := os.Getenv("SERVER_PORT")
	jwtSecret := os.Getenv("JWT_SECRET")

	if dbUrl == "" {
		return &AppConfig{}, errors.New("Failed to load DATABASE_URL from ENV")
	}

	if port == "" {
		return &AppConfig{}, errors.New("Failed to load PORT ENV")
	}

	if jwtSecret == "" {
		return &AppConfig{}, errors.New("Failed to load JWT_SECRET ENV")
	}

	return &AppConfig{
		DatabaseUrl: dbUrl,
		Port:        port,
	}, nil
}
