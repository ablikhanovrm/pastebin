package http

import (
	"errors"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	"github.com/ablikhanovrm/pastebin/internal/models/user"
	pasteService "github.com/ablikhanovrm/pastebin/internal/service/paste"
)

func MapError(err error) (int, string) {
	switch {
	// ---- paste domain ----
	case errors.Is(err, paste.ErrNotFound):
		return 404, "paste not found"

	case errors.Is(err, paste.ErrForbidden):
		return 403, "forbidden"

	case errors.Is(err, paste.ErrExpired):
		return 410, "paste expired"
	case errors.Is(err, pasteService.ErrUpdate):
		return 400, "failed update paste"

	// ---- user domain ----
	case errors.Is(err, user.ErrNotFound):
		return 404, "user not found"

	case errors.Is(err, pasteService.ErrUploadFailed):
		return 400, "failed to upload paste content"

	default:
		return 500, "internal error"
	}
}
