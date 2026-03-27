package auth

import (
	store "github.com/arun-kumar21/koffee/internal/store/sqlc/gen"
)

type Service struct {
	queries *store.Queries
	tokens *TokenManager
}


func NewService(queries *store.Queries, tokens *TokenManager) *Service {
	return &Service{queries: queries, tokens: tokens}
}


