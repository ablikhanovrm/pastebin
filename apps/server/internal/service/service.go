package service

import (
	"github.com/ablikhanovrm/pastebin/internal/repository"
	"github.com/ablikhanovrm/pastebin/internal/service/auth"
	"github.com/ablikhanovrm/pastebin/internal/service/paste"
	"github.com/ablikhanovrm/pastebin/internal/service/user"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
)

type Services struct {
	Auth  *auth.AuthService
	User  *user.UserService
	Paste *paste.PasteService
}

func NewService(
	repo *repository.Repository,
	jwtManager *jwt.Manager,
) *Services {
	return &Services{
		Auth:  auth.NewAuthService(repo.User, repo.Auth, jwtManager),
		User:  user.NewUserService(repo.User),
		Paste: paste.NewPasteService(repo.Paste),
	}
}
