package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)


type Handler struct {
	service *Service
	tokens *TokenManager
}

func NewHandler (service *Service, tokens *TokenManager) *Handler {
	return &Handler{service: service, tokens: tokens}
}

func MountRoutes (r chi.Router, handler *Handler) {
	r.Route("/api/v1/auth", func (r chi.Router) {
		r.Post("/register", handler.handleRegister)
		r.Post("/login", handler.handleLogin)
		//TODO: middleware
	})
}


func (h *Handler) handleRegister (w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleLogin (w http.ResponseWriter, r *http.Request) {

}
