package paste

import "errors"

var (
	ErrUploadFailed = errors.New("failed to upload paste content")
	ErrUpdate       = errors.New("failed to update paste")
)
