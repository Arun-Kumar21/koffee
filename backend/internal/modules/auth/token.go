package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidToken = errors.New("Invaild Token")
	ErrExpiredToken = errors.New("Token expire")
)

type TokenClaims struct {
	UserId string `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret          []byte
	refreshTokenTTL time.Duration
	accessTokenTTL  time.Duration
}

func NewTokenManager(secret string, refreshTokenTTL time.Duration, accessTokenTTL time.Duration) *TokenManager {
	return &TokenManager{
		secret:          []byte(secret),
		refreshTokenTTL: refreshTokenTTL,
		accessTokenTTL:  accessTokenTTL,
	}
}

func (tm *TokenManager) GenerateAccessToken(userId uuid.UUID) (string, error) {
	claims := TokenClaims{
		UserId: userId.String(),
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secret)
}

func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (tm *TokenManager) HashRefreshToken(token string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(token), bcryptCost)
	return string(hash)
}

func (tm *TokenManager) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return tm.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Type != "access" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
