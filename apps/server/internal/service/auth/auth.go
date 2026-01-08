package auth

import (
	"github.com/ablikhanovrm/pastebin/internal/repository/auth"
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
)

type Service interface {
	Login(email, password string) (*Tokens, error)
	Refresh(refreshToken string) (*Tokens, error)
	Logout(refreshToken string) error
}

type AuthService struct {
	users    user.UserRepository
	tokens   *jwt.Manager
	authRepo auth.AuthRepository
}

func NewAuthService(
	users user.UserRepository,
	authRepo auth.AuthRepository,
	tokens *jwt.Manager,
) *AuthService {
	return &AuthService{
		users:    users,
		authRepo: authRepo,
		tokens:   tokens,
	}
}

func (s *AuthService) Login(email string, password string) (*Tokens, error) {
}

func (s *AuthService) Refresh(refreshToken string) (*Tokens, error) {
}

func (s *AuthService) Logout(refreshToken string) error {
}
