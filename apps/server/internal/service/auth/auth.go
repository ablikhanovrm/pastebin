package auth

import (
	"context"
	"time"

	"github.com/ablikhanovrm/pastebin/internal/repository/auth"
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/ablikhanovrm/pastebin/pkg/random"
	"github.com/ablikhanovrm/pastebin/pkg/security"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
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

func (s *AuthService) Login(ctx context.Context, email string, password string) (*Tokens, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nil, ErrUserAlreadyExists
	}

	token, err := s.tokens.Generate(user.ID, time.Second*60)
	if err != nil {
		return nil, err
	}

	if !security.CheckPassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	refreshToken, err := random.GenerateRefreshToken(32)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
}
