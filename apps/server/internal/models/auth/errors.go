package auth

import "errors"

var (
	ErrTokenExpired   = errors.New("token expired")
	ErrReauthRequired = errors.New("session expired")
)
