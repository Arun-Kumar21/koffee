package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)


type Handler struct {
	service *Service
	tokens *TokenManager
}

const maxAuthBodyBytes = 4 * 1024;

type UserResponse struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarUrl string `json:"avatar_url"`
	Role      string `json:"role"`
}

type registerRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
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
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	var parsedRequest registerRequest
	if err := json.NewDecoder(r.Body).Decode(&parsedRequest); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error" : "Request validation failed"})
		return
	}


}

func (h *Handler) handleLogin (w http.ResponseWriter, r *http.Request) {

}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

