package repository

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/repository/paste"
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
)

type Repository struct {
	user  *user.UserRepository
	paste *paste.PasteRepository
}

func NewRepository(q *dbgen.Queries) *Repository {
	return &Repository{
		user:  user.NewUserRepository(q),
		paste: paste.NewPasteRepository(q),
	}
}
