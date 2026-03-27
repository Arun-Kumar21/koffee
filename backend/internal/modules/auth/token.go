package auth

import "time"

type TokenManager struct {
	secret []byte
	refreshTokenTTL time.Duration
	accessTokenTTL time.Duration
}

func NewTokenManager (secret string, refreshTokenTTL time.Duration, accessTokenTTL time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		refreshTokenTTL: refreshTokenTTL,
		accessTokenTTL: accessTokenTTL,
	}
}

