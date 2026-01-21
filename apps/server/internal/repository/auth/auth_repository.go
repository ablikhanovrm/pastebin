package auth

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, userID int64, token models.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, token string) error
	IsRefreshTokenValid(ctx context.Context, token string) (bool, error)
}

type SqlcAuthRepository struct {
	q *dbgen.Queries
}

func NewSqlcAuthRepository(q *dbgen.Queries) *SqlcAuthRepository {
	return &SqlcAuthRepository{q: q}
}

func (r *SqlcAuthRepository) SaveRefreshToken(ctx context.Context, userID int64, token models.RefreshToken) (int64, error) {
	params := dbgen.CreateRefreshTokenParams{
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		UserAgent: toPgText(token.UserAgent),
		IpAddress: toNetIp(token.IPAddress),
		ExpiresAt: pgtype.Timestamptz{
			Time:  token.ExpiresAt,
			Valid: true,
		},
	}

	row, err := r.q.CreateRefreshToken(ctx, params)
	if err != nil {
		return 0, err
	}

	return row.ID, nil
}

func (r *SqlcAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return nil
}

func (r *SqlcAuthRepository) IsRefreshTokenValid(ctx context.Context, token string) (bool, error) {
	return false, nil
}
