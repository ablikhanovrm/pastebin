package user

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/user"
)

func mapUserByEmail(u dbgen.User) *user.User {
	return &user.User{
		Id:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt.Time,
		PasswordHash: u.PasswordHash,
	}
}

func mapUserById(u dbgen.User) *user.User {
	return &user.User{
		Id:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		CreatedAt:    u.CreatedAt.Time,
		PasswordHash: u.PasswordHash,
	}
}
