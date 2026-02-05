package user

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/user"
	userrepo "github.com/ablikhanovrm/pastebin/internal/repository/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type UserService interface {
	Create(ctx context.Context, u *user.User) (*user.User, error)
	FindById(ctx context.Context, id int64) (*user.User, error)
}

type Service struct {
	log zerolog.Logger
	db  *pgxpool.Pool
}

func (s *Service) repo(db dbgen.DBTX) *userrepo.SqlcUserRepository {
	return userrepo.NewSqlcUserRepository(db, s.log)
}

func NewUserService(db *pgxpool.Pool, log zerolog.Logger) *Service {
	return &Service{db: db, log: log}
}
