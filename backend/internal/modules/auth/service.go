package auth

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"

	store "github.com/arun-kumar21/koffee/internal/store/sqlc/gen"
	
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

const bcryptCost = 12
const maxPasswordBytes = 72

type Service struct {
	queries *store.Queries
	tokens *TokenManager
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	AccessToken string
	RefreshToken string
	User  store.User
}

func NewService(queries *store.Queries, tokens *TokenManager) *Service {
	return &Service{queries: queries, tokens: tokens}
}


func (s *Service) Register (ctx context.Context, input RegisterInput) (AuthResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return AuthResult{}, errors.New("Email is required")
	}

	if !emailRegex.MatchString(email) {
		return AuthResult{}, errors.New("Invalid email format")
	}

	if len(input.Password) < 8 {
		return AuthResult{}, errors.New("Password must be at least 8 characters")
	}

	if len([]byte(input.Password)) > maxPasswordBytes {
		return AuthResult{}, errors.New("Password must not exceed 72 bytes")
	}

	if len(input.Name) > 30 {
		return AuthResult{}, errors.New("Name must not be greater than 30 characters")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		return AuthResult{}, err
	}

	created, err := s.queries.CreateUser(ctx, store.CreateUserParams{
		Email: input.Email,
		Password: hashPassword,
		Name: input.Name,
		AvatarUrl: stringToNullString(input.AvatarUrl),
	})
	
	if err != nil {
		// pq error code 23505 = unique_violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return AuthResult{}, errors.New("email already registered")
		}
		return AuthResult{}, err
	}

	//TODO: Return Token after login
}




func stringToNullString(value string) sql.NullString {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}
