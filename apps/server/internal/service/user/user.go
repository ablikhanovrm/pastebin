package user

import (
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
	"github.com/rs/zerolog"
)

type UserService struct {
	repo   user.UserRepository
	logger zerolog.Logger
}

func NewUserService(repo user.UserRepository, logger zerolog.Logger) *UserService {
	return &UserService{repo: repo, logger: logger}
}
