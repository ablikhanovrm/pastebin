package user

import (
	"context"
	"errors"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/user"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id int64) (*user.User, error)
	Create(ctx context.Context, u *user.User) error
}

type SqlcUserRepository struct {
	q      *dbgen.Queries
	logger zerolog.Logger
}

func NewSqlcUserRepository(db dbgen.DBTX, logger zerolog.Logger) *SqlcUserRepository {
	return &SqlcUserRepository{
		q:      dbgen.New(db),
		logger: logger,
	}
}
func (r *SqlcUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	foundUser, err := r.q.GetUserByEmail(ctx, email)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return mapUserByEmail(foundUser), nil
}

func (r *SqlcUserRepository) FindByID(ctx context.Context, id int64) (*user.User, error) {
	foundUser, err := r.q.GetUserById(ctx, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return mapUserById(foundUser), nil
}

func (r *SqlcUserRepository) Create(ctx context.Context, u *user.User) error {

}
