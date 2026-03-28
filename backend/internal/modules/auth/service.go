package auth

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	store "github.com/arun-kumar21/koffee/internal/store/sqlc/gen"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

const bcryptCost = 12
const maxPasswordBytes = 72

type Service struct {
	queries *store.Queries
	tokens  *TokenManager
}

type RegisterInput struct {
	Name      string
	Email     string
	Password  string
	AvatarUrl string
}

type RegisterResult struct {
	UserId string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthUserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	AvatarUrl sql.NullString `json:"avatar_url"`
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         AuthUserResponse 
}

func NewService(queries *store.Queries, tokens *TokenManager) *Service {
	return &Service{queries: queries, tokens: tokens}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (RegisterResult, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" {
		return RegisterResult{}, errors.New("Email is required")
	}

	if !emailRegex.MatchString(email) {
		return RegisterResult{}, errors.New("Invalid email format")
	}

	if len(input.Password) < 8 {
		return RegisterResult{}, errors.New("Password must be at least 8 characters")
	}

	if len([]byte(input.Password)) > maxPasswordBytes {
		return RegisterResult{}, errors.New("Password must not exceed 72 bytes")
	}

	if len(input.Name) > 30 {
		return RegisterResult{}, errors.New("Name must not be greater than 30 characters")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		return RegisterResult{}, err
	}

	created, err := s.queries.CreateUser(ctx, store.CreateUserParams{
		Email:     input.Email,
		Password:  string(hashPassword),
		Name:      input.Name,
		AvatarUrl: stringToNullString(input.AvatarUrl),
	})

	if err != nil {
		// pq error code 23505 = unique_violation
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return RegisterResult{}, errors.New("email already registered")
		}
		return RegisterResult{}, err
	}

	return RegisterResult{UserId: created.ID.String()}, nil

}

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	user, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return AuthResult{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return AuthResult{}, errors.New("invalid credentials")
	}

	accessToken, err := s.tokens.GenerateAccessToken(user.ID)
	if err != nil {
		return AuthResult{}, err
	}

	refreshToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return AuthResult{}, err
	}

	tokenHash := s.tokens.HashRefreshToken(refreshToken)
	expiresAt := time.Now().Add(s.tokens.refreshTokenTTL)

	_, refTokenErr := s.queries.CreateRefreshToken(ctx, store.CreateRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})

	if refTokenErr != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		User:         AuthUserResponse{
			ID: user.ID,
			Email: user.Email,
			Name: user.Name,
			AvatarUrl: user.AvatarUrl,
		},
	}, nil
}

func stringToNullString(value string) sql.NullString {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: trimmed, Valid: true}
}
