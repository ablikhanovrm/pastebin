package auth

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string
	Revoked   bool

	UserAgent *string
	IPAddress *string

	ExpiresAt        time.Time
	CreatedAt        time.Time
	SessionExpiresAt time.Time
}
