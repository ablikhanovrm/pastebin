package handler

import (
	"time"

	"github.com/ablikhanovrm/pastebin/internal/models/paste"
)

type CreatePasteRequest struct {
	Title      string           `json:"title"`
	Content    string           `json:"content"`
	Syntax     paste.Syntax     `json:"syntax"`
	Visibility paste.Visibility `json:"visibility"`
	MaxViews   *int32           `json:"maxViews,omitempty"`
	ExpireAt   *time.Time       `json:"expireAt,omitempty"`
}

type UpdatePasteRequest struct {
	Title      string           `json:"title"`
	Content    string           `json:"content"`
	Syntax     paste.Syntax     `json:"syntax"`
	Visibility paste.Visibility `json:"visibility"`
	MaxViews   *int32           `json:"maxViews,omitempty"`
	ExpireAt   *time.Time       `json:"expireAt,omitempty"`
}
