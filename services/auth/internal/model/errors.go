package model

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
)
