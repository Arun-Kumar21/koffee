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
}

type registerRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarUrl string `json:"avatar_url"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}



func NewHandler (service *Service, tokens *TokenManager) *Handler {
	return &Handler{service: service, tokens: tokens}
}

func MountRoutes (r chi.Router, handler *Handler) {
	r.Route("/api/v1/auth", func (r chi.Router) {
		r.Post("/register", handler.handleRegister)
		r.Post("/login", handler.handleLogin)
	})
}


func (h *Handler) handleRegister (w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	var parsedRequest registerRequest
	if err := json.NewDecoder(r.Body).Decode(&parsedRequest); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error" : "Request validation failed"})
		return
	}
	
	user, err := h.service.Register(r.Context(), RegisterInput(parsedRequest))
	if err != nil {
		switch err.Error() {
		case "email is required", "invalid email format",
		     "password must be at least 8 characters", "password must not exceed 72 bytes":
			 writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case "email already registered":
			writeJSON(w, http.StatusConflict, map[string]string{"error": "An account with this email already exists"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Something went wrong. Please try again."})
		}

		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"userId": user.UserID})
}

func (h *Handler) handleLogin (w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxAuthBodyBytes)
	var parsedRequest LoginInput
	if err := json.NewDecoder(r.Body).Decode(&parsedRequest); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string {
			"error": "Request validation failed",
		})
	}

	authResult, err := h.service.Login(r.Context(), parsedRequest)
	if err != nil {
		switch err.Error() {
			case "invalid credentials":
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			default:
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}

		return 
	}

	writeJSON(w, http.StatusOK, authResponse{
		RefreshToken: authResult.RefreshToken,
		AccessToken: authResult.AccessToken,
		User: authResult.User,
	})
}

func (h *Handler) handleRefresh (w http.ResponseWriter, r *http.Request) {
 	// TODO: To implement
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

