package config

import "os"




type Config struct {
	ServerPort string
	DatabaseUrl string
	JWTSecret string
}


func Load() *Config {
	serverPort := getEnv("SERVER_PORT", "8080")
	dbUrl := getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/auth?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "SECRET")

	return &Config{
		ServerPort: serverPort,
		DatabaseUrl: dbUrl,
		JWTSecret: jwtSecret,
	}
}


func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
