package auth

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
)

type AuthRepository interface {
	SaveRefreshToken(ctx context.Context, userID int64, token string) error
	DeleteRefreshToken(ctx context.Context, token string) error
	IsRefreshTokenValid(ctx context.Context, token string) (bool, error)
}

type SqlcAuthRepository struct {
	q *dbgen.Queries
}

func NewSqlcAuthRepository(q *dbgen.Queries) *SqlcAuthRepository {
	return &SqlcAuthRepository{q: q}
}

func (*SqlcAuthRepository) SaveRefreshToken(ctx context.Context, userID int64, token string) error {
	return nil
}

func (*SqlcAuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	return nil
}

func (*SqlcAuthRepository) IsRefreshTokenValid(ctx context.Context, token string) (bool, error) {
	return false, nil
}
