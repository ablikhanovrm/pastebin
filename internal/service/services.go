package service

import (
	"github.com/ablikhanovrm/pastebin/internal/repository"
	"github.com/ablikhanovrm/pastebin/internal/repository/cache"
	"github.com/ablikhanovrm/pastebin/internal/service/auth"
	"github.com/ablikhanovrm/pastebin/internal/service/paste"
	"github.com/ablikhanovrm/pastebin/internal/service/storage"
	"github.com/ablikhanovrm/pastebin/internal/service/user"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Services struct {
	Auth  *auth.Service
	User  *user.Service
	Paste *paste.Service
}

func NewServices(
	repo *repository.Repository,
	jwtManager *jwt.Manager,
	db *pgxpool.Pool,
	s3Storage *storage.Service,
	cache *cache.RedisCache,
) *Services {
	return &Services{
		Auth:  auth.NewAuthService(repo.User, jwtManager, db, cache),
		User:  user.NewUserService(db, cache),
		Paste: paste.NewPasteService(db, s3Storage, cache),
	}
}
