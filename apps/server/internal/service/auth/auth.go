package auth

import (
	"context"
	"time"

	"github.com/ablikhanovrm/pastebin/internal/repository/auth"
	"github.com/ablikhanovrm/pastebin/internal/repository/user"
	"github.com/ablikhanovrm/pastebin/pkg/hash"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/ablikhanovrm/pastebin/pkg/random"
	"github.com/ablikhanovrm/pastebin/pkg/security"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthService struct {
	users  user.UserRepository
	repo   auth.AuthRepository
	tokens *jwt.Manager
	db     *pgxpool.Pool
}

func NewAuthService(
	users user.UserRepository,
	repo auth.AuthRepository,
	tokens *jwt.Manager,
	db *pgxpool.Pool,
) *AuthService {
	return &AuthService{
		repo:   repo,
		users:  users,
		tokens: tokens,
		db:     db,
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
	hash := hash.HashRefreshToken(refreshToken)
	err := s.repo.RevokeRefreshTokenByHash(ctx, hash)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*Tokens, error) {
	hash := hash.HashRefreshToken(refreshToken)
	rt, err := s.repo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	if time.Now().After(rt.SessionExpiresAt) {
		return nil, ErrReauthRequired
	}

	// TRANSACTION
	// ├─ revoke OLD refresh
	// ├─ create NEW refresh
	// └─ generate NEW access
}
