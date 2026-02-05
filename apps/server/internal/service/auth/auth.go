package auth

import (
	"context"
	"errors"
	"time"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/auth"
	"github.com/ablikhanovrm/pastebin/internal/models/user"
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

type AuthService interface {
	Login(ctx context.Context, email, password string) (*Tokens, error)
	Register(ctx context.Context, input RegisterInput) (*Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*Tokens, error)
	Logout(ctx context.Context, refreshToken string) error
}

type Service struct {
	users  userrepo.UserRepository
	tokens *jwt.Manager
	db     *pgxpool.Pool
	log    zerolog.Logger
}

func NewAuthService(
	users userrepo.UserRepository,
	tokens *jwt.Manager,
	db *pgxpool.Pool,
	log zerolog.Logger,
) *Service {
	return &Service{
		users:  users,
		tokens: tokens,
		db:     db,
		log:    log,
	}
}

// repo helper
func (s *Service) repo(db dbgen.DBTX) *authrepo.SqlcAuthRepository {
	return authrepo.NewSqlcAuthRepository(db, s.log)
}

func (s *Service) Login(ctx context.Context, params LoginInput) (*Tokens, error) {
	foundUser, err := s.users.FindByEmail(ctx, params.Email)
	if err != nil {
		return nil, err
	}

	if !security.CheckPassword(foundUser.PasswordHash, params.Password) {
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

	_, err = s.repo(s.db).CreateRefreshToken(ctx, auth.RefreshToken{
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 15),
		SessionExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		TokenHash:        refreshToken,
		UserAgent:        &params.UserAgent,
		UserID:           foundUser.Id,
		IPAddress:        &params.IP,
	})

	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	hashRt := hash.HashRefreshToken(refreshToken)
	err := s.repo(s.db).RevokeRefreshTokenByHash(ctx, hashRt)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string, ip string, ua string) (*Tokens, error) {
	hashRt := hash.HashRefreshToken(refreshToken)
	rt, err := s.repo(s.db).GetRefreshTokenByHash(ctx, hashRt)
	if err != nil {
		return nil, err
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, auth.ErrTokenExpired
	}
	if time.Now().After(rt.SessionExpiresAt) {
		return nil, auth.ErrReauthRequired
	}

	foundUser, err := s.users.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, err
	}

	newRtHash, err := random.GenerateRefreshToken(32)
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

	if err := s.repo(tx).RevokeRefreshTokenByHash(ctx, hashRt); err != nil {
		return nil, err
	}

	_, err = s.repo(tx).CreateRefreshToken(ctx, auth.RefreshToken{
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
		s.log.Warn().Err(err).
			Int64("user_id", foundUser.Id).
			Msg("failed to generate access token")

		token = ""
	}

	return &Tokens{AccessToken: token, RefreshToken: newRtHash}, nil
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*Tokens, error) {
	hashPass, err := security.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	userId, err := s.users.Create(ctx, user.User{
		Email:        input.Email,
		PasswordHash: hashPass,
		Name:         input.Name,
	})

	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			return nil, user.ErrUserAlreadyExists
		}
		return nil, err
	}

	access, err := s.tokens.Generate(userId, 15*time.Minute)
	if err != nil {
		s.log.Warn().Err(err).Msg("failed generate access")
	}

	refresh, err := random.GenerateRefreshToken(32)
	if err != nil {
		return nil, err
	}

	_, err = s.repo(s.db).CreateRefreshToken(ctx, auth.RefreshToken{
		ExpiresAt:        time.Now().Add(time.Hour * 24 * 30),
		SessionExpiresAt: time.Now().Add(time.Hour * 24 * 30),
		TokenHash:        refresh,
		UserAgent:        &input.UserAgent,
		UserID:           userId,
		IPAddress:        &input.IP,
	})

	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
