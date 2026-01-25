package auth

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenExpired       = errors.New("token expired")
	ErrReauthRequired     = errors.New("session expired")
	ErrRefreshExpired     = errors.New("refresh token expired")
)
