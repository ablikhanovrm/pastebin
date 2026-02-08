package paste

import (
	"time"

	"github.com/google/uuid"
)

type Syntax string

const (
	SyntaxPlain Syntax = "plain"
	SyntaxCode  Syntax = "code"
)

type Visibility string

const (
	VisibilityPublic   Visibility = "public"
	VisibilityUnlisted Visibility = "unlisted"
	VisibilityPrivate  Visibility = "private"
)

type PasteStatus string

const (
	Active     PasteStatus = "active"
	Terminated PasteStatus = "terminated"
)

type Paste struct {
	Id     int64
	Uuid   uuid.UUID
	UserId int64

	Title string
	S3Key string

	Syntax     Syntax
	Visibility Visibility

	ViewsCount int32
	MaxViews   *int32 // null = бесконечно

	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
