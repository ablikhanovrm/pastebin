package user

import dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"

type UserRepository struct {
	q *dbgen.Queries
}

func NewUserRepository(q *dbgen.Queries) *UserRepository {
	return &UserRepository{q: q}
}
