package user

import (
	"context"
	"errors"
	"time"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, id int64) (*user.User, error)
	Create(ctx context.Context, u user.User) (int64, error) // TODO: use custom input for user
}

type SqlcUserRepository struct {
	q   *dbgen.Queries
	log zerolog.Logger
}

func NewSqlcUserRepository(db dbgen.DBTX, log zerolog.Logger) *SqlcUserRepository {
	return &SqlcUserRepository{
		q:   dbgen.New(db),
		log: log,
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

func (r *SqlcUserRepository) Create(ctx context.Context, u user.User) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	params := dbgen.CreateUserParams{
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}

	row, err := r.q.CreateUser(ctx, params)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return 0, user.ErrUserAlreadyExists
			}
		}

		return 0, err
	}

	return row.ID, nil
}
