package repository

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/repository/auth"
	"github.com/ablikhanovrm/pastebin/internal/repository/paste"
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
	"github.com/rs/zerolog"
)

type Repository struct {
	User  user.UserRepository
	Paste paste.PasteRepository
	Auth  auth.AuthRepository
}

func NewRepository(db dbgen.DBTX, logger zerolog.Logger) *Repository {
	return &Repository{
		User:  user.NewSqlcUserRepository(db, logger),
		Paste: paste.NewSqlcPasteRepository(db, logger),
		Auth:  auth.NewSqlcAuthRepository(db, logger),
	}
}
