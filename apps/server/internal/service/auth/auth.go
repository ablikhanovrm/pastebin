package auth

import (
	"context"
	"time"

	auth "github.com/ablikhanovrm/pastebin/internal/models/auth"
	user "github.com/ablikhanovrm/pastebin/internal/models/user"
	authrepo "github.com/ablikhanovrm/pastebin/internal/repository/auth"
	userrepo "github.com/ablikhanovrm/pastebin/internal/repository/user"
	"github.com/ablikhanovrm/pastebin/pkg/hash"
	"github.com/ablikhanovrm/pastebin/pkg/jwt"
	"github.com/ablikhanovrm/pastebin/pkg/random"
	"github.com/ablikhanovrm/pastebin/pkg/security"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Service interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
}

type AuthService struct {
	users  userrepo.UserRepository
	repo   authrepo.AuthRepository
	tokens *jwt.Manager
	db     *pgxpool.Pool
	logger zerolog.Logger
}

func NewAuthService(
	users userrepo.UserRepository,
	repo authrepo.AuthRepository,
	tokens *jwt.Manager,
	db *pgxpool.Pool,
	logger zerolog.Logger,
) *AuthService {
	return &AuthService{
		repo:   repo,
		users:  users,
		tokens: tokens,
		db:     db,
		logger: logger,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*Tokens, error) {
	foundUser, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !security.CheckPassword(foundUser.PasswordHash, password) {
		return nil, user.ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(foundUser.Id, time.Minute*15)
	if err != nil {
		return nil, err
	}

	refreshToken, err := random.GenerateRefreshToken(32)
	if err != nil {
		return nil, err
	}
	//TODO: save rt to db

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

func (s *AuthService) Refresh(ctx context.Context, refreshToken string, ip string, ua string) (*Tokens, error) {
	hash := hash.HashRefreshToken(refreshToken)
	rt, err := s.repo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, auth.ErrTokenExpired
	}
	if time.Now().After(rt.SessionExpiresAt) {
		return nil, auth.ErrReauthRequired
	}

	foundUser, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	})

	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	s.repo.RevokeRefreshTokenByHash(ctx, hash)

	newRtHash, err := random.GenerateRefreshToken(32)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.CreateRefreshToken(ctx, rt.UserID, auth.RefreshToken{
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 30),
		SessionExpiresAt: rt.SessionExpiresAt,
		TokenHash:        newRtHash,
		UserAgent:        &ua,
		UserID:           rt.UserID,
		IPAddress:        &ip,
	})

	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, err
	}

	token, err := s.tokens.Generate(foundUser.Id, time.Minute*15)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to generate access token")
	}

	return &Tokens{AccessToken: token, RefreshToken: newRtHash}, nil
}
