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

func NewRepository(db dbgen.DBTX) *Repository {
	return &Repository{
		User:  user.NewSqlcUserRepository(db),
		Paste: paste.NewSqlcPasteRepository(db),
		Auth:  auth.NewSqlcAuthRepository(db),
	}
}
