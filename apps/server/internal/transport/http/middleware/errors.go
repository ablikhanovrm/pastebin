package middleware

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)
