package models

import "time"

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

type UserStatus string

const (
	Active     UserStatus = "active"
	Terminated UserStatus = "Terminated"
)

type Paste struct {
	Id           int64
	PasteUrl     string
	UserId       int64
	Title        string
	Content      string
	Syntax       Syntax
	Visibility   Visibility
	PasswordHash *string
	MaxViews     *int32
	ViewsCount   int32
	ExpiresAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Status       UserStatus
}
