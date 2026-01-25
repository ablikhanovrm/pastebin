package auth

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, userID int64, token models.RefreshToken) (int64, error)
	RevokeRefreshTokenByHash(ctx context.Context, token string) error
	GetRefreshTokenByHash(ctx context.Context, token string) (*models.RefreshToken, error)
}

type SqlcAuthRepository struct {
	q *dbgen.Queries
}

func NewSqlcAuthRepository(db dbgen.DBTX) *SqlcAuthRepository {
	return &SqlcAuthRepository{
		q: dbgen.New(db),
	}
}

func (r *SqlcAuthRepository) CreateRefreshToken(ctx context.Context, userID int64, token models.RefreshToken) (int64, error) {
	params := dbgen.CreateRefreshTokenParams{
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		UserAgent: toPgText(token.UserAgent),
		IpAddress: toNetIp(token.IPAddress),
		ExpiresAt: pgtype.Timestamptz{ // TTL
			Time:  token.ExpiresAt,
			Valid: true,
		},
		SessionExpiresAt: pgtype.Timestamptz{ // absolute TTL
			Time:  token.SessionExpiresAt,
			Valid: true,
		},
	}

	row, err := r.q.CreateRefreshToken(ctx, params)
	if err != nil {
		return 0, err
	}

	return row.ID, nil
}

func (r *SqlcAuthRepository) RevokeRefreshTokenByHash(ctx context.Context, token string) error {
	err := r.q.RevokeRefreshToken(ctx, token)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlcAuthRepository) GetRefreshTokenByHash(ctx context.Context, token string) (*models.RefreshToken, error) {
	refreshToken, err := r.q.GetRefreshTokenByHash(ctx, token)
	if err != nil {
		return nil, err
	}

	return mapRefreshToken(refreshToken), nil
}
