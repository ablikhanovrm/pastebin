package repository

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/repository/auth"
	"github.com/ablikhanovrm/pastebin/internal/repository/paste"
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
)

type Repository struct {
	User  user.UserRepository
	Paste paste.PasteRepository
	Auth  auth.AuthRepository
}

func NewRepository(q *dbgen.Queries) *Repository {
	return &Repository{
		User:  user.NewSqlcUserRepository(q),
		Paste: paste.NewSqlcPasteRepository(q),
		Auth:  auth.NewSqlcAuthRepository(q),
	}
}
