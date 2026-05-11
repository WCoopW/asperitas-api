package auth

import "errors"

var (
	ErrMissingSecret  = errors.New("JWT_SECRET is not set")
	ErrSignFailed     = errors.New("failed to sign token")
	ErrInvalidToken   = errors.New("invalid jwt token")
	ErrExpiredToken   = errors.New("expired jwt token")
	ErrNotInitialized = errors.New("jwt is not initialized")
	ErrInvalidUserID  = errors.New("user id is required")
	ErrInvalidTTL     = errors.New("jwt expire must be greater than zero")
)
