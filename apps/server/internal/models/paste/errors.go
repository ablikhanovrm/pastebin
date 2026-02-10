package paste

import "errors"

var (
	ErrNotFound  = errors.New("paste not found")
	ErrForbidden = errors.New("forbidden")
	ErrInvalidID = errors.New("invalid id")
)
