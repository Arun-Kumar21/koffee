package main

import (
	"log"

	"github.com/Arun-Kumar21/koffee/services/auth/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load ENV: %v", err)
	}

	cfg := config.Load()

}
