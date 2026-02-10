package handler

import "errors"

var (
	ErrMissingIDParam = errors.New("missing id query param")
	ErrInvalidQuery   = errors.New("invalid query param")
	ErrInvalidJSON    = errors.New("invalid json body")
)
