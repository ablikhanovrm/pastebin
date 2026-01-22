package auth

import (
	"net/netip"
	"time"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func toNetIp(s *string) *netip.Addr {
	if s == nil {
		return nil
	}

	addr, err := netip.ParseAddr(*s)
	if err != nil {
		return nil
	}

	return &addr
}

func toPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}

	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

func mapRefreshToken(row dbgen.GetRefreshTokenByHashRow) *models.RefreshToken {
	var ua *string
	if row.UserAgent.Valid {
		ua = &row.UserAgent.String
	}

	var expiresAt time.Time
	if row.ExpiresAt.Valid {
		expiresAt = row.ExpiresAt.Time
	}

	var ip string

	if row.IpAddress.IsValid() {
		ip = row.IpAddress.String()
	}

	return &models.RefreshToken{
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		UserAgent: ua,
		IPAddress: &ip,
		ExpiresAt: expiresAt,
	}
}
