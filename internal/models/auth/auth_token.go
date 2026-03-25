package auth

import "time"

type RefreshToken struct {
	ID        int64
	UserID    int64
	TokenHash string

	UserAgent *string
	IPAddress *string

	RevokedAt        time.Time
	ExpiresAt        time.Time
	CreatedAt        time.Time
	SessionExpiresAt time.Time
}
