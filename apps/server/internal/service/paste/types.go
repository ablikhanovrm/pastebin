package paste

import (
	"time"
)

type UpdatePasteInput struct {
	Title      string
	Syntax     string
	Visibility string
	MaxViews   *int32
	ExpireAt   *time.Time
}
