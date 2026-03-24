package model

import "time"

type User struct {
	ID        string  `json:"id"`
	Name      string `json:"name"`
	Email     string  `json:"email"`
	Password  string  `json:"-"`
	Avatar_url string `json:"avatar_url"`

	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Avatar_url string `json:"avatar_url"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
