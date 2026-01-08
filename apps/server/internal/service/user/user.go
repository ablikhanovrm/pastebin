package user

import (
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
)

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{repo: repo}
}
