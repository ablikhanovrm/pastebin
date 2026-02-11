package paste

import "errors"

var (
	ErrNotFound  = errors.New("paste not found")
	ErrForbidden = errors.New("forbidden")
	ErrExpired   = errors.New("paste expired")
)
