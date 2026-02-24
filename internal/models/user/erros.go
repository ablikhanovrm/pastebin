package user

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrNotFound           = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
