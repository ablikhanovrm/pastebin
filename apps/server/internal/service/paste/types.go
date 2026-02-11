package paste

import (
	"time"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
)

type UpdatePasteInput struct {
	Title      string
	Syntax     paste.Syntax
	Visibility paste.Visibility
	MaxViews   *int32
	ExpireAt   *time.Time
}

type CreatePasteInput struct {
	Title      string
	Content    string
	Syntax     paste.Syntax
	Visibility paste.Visibility
	MaxViews   *int32
	ExpireAt   *time.Time
}
