package main

import (
	"log"

	"github.com/arun-kumar21/koffee/config"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	_, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}

	log.Printf("Server started successfully")
}
