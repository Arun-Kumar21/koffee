package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/arun-kumar21/koffee/config"
	"github.com/arun-kumar21/koffee/internal/modules/auth"
	store "github.com/arun-kumar21/koffee/internal/store/sqlc/gen"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	conn, err := sql.Open("postgres", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect databse: %v", err)
	}

	queries := store.New(conn)
	tokenManger := auth.NewTokenManager(cfg.JWTSecret, 30 * time.Minute, 7 * 24 * time.Hour)
	authService := auth.NewService(queries, tokenManger)
	authHandler := auth.NewHandler(authService, tokenManger)


	r := chi.NewRouter()
	r.Get("/health", func (w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if err := conn.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("db_unavailable"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})


	auth.MountRoutes(r, authHandler)

	log.Printf("Server started successfully")
}
