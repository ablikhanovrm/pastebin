package service

import (
	"github.com/ablikhanovrm/pastebin/internal/repository"
	"github.com/ablikhanovrm/pastebin/internal/service/auth"
	"github.com/ablikhanovrm/pastebin/internal/service/paste"
	"github.com/ablikhanovrm/pastebin/internal/service/user"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Services struct {
	Auth  *auth.AuthService
	User  *user.UserService
	Paste *paste.PasteService
}

func NewServices(
	repo *repository.Repository,
	jwtManager *jwt.Manager,
	db *pgxpool.Pool,
	logger zerolog.Logger,
) *Services {
	authLogger := logger.With().
		Str("service", "auth").
		Logger()

	userLogger := logger.With().
		Str("service", "user").
		Logger()

	pasteLogger := logger.With().
		Str("service", "paste").
		Logger()

	return &Services{
		Auth:  auth.NewAuthService(repo.User, repo.Auth, jwtManager, db, authLogger),
		User:  user.NewUserService(repo.User, userLogger),
		Paste: paste.NewPasteService(repo.Paste, pasteLogger),
	}
}
