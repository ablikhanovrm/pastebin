package user

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*dbgen.User, error)
	FindByID(ctx context.Context, id int64) (*dbgen.User, error)
	Create(ctx context.Context, u *dbgen.User) error
}

type SqlcUserRepository struct {
	db dbgen.DBTX
}

func NewSqlcUserRepository(db dbgen.DBTX) *SqlcUserRepository {
	return &SqlcUserRepository{db: db}
}

func (*SqlcUserRepository) FindByEmail(ctx context.Context, email string) (*dbgen.User, error) {
	return nil, nil
}

func (*SqlcUserRepository) FindByID(ctx context.Context, id int64) (*dbgen.User, error) {
	return nil, nil
}

func (*SqlcUserRepository) Create(ctx context.Context, u *dbgen.User) error {
	return nil
}
