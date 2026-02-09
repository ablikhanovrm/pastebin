package paste

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	pasterepo "github.com/ablikhanovrm/pastebin/internal/repository/paste"
	"github.com/ablikhanovrm/pastebin/internal/service/storage"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type PasteService interface {
	Create(ctx context.Context, u *paste.Paste) (*paste.Paste, error)
	GetByID(ctx context.Context, id uuid.UUID) (*paste.Paste, error)
	GetPastes(ctx context.Context, userId int64) ([]*paste.Paste, error)
	GetMyPastes(ctx context.Context, userId int64) ([]*paste.Paste, error)
	Update(ctx context.Context, paste *paste.Paste) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	db        *pgxpool.Pool
	s3Storage *storage.Service
	log       zerolog.Logger
}

func NewPasteService(db *pgxpool.Pool, s3Storage *storage.Service, log zerolog.Logger) *Service {
	return &Service{db: db, s3Storage: s3Storage, log: log}
}

// repo helper
func (s *Service) repo(db dbgen.DBTX) *pasterepo.SqlcPasteRepository {
	return pasterepo.NewSqlcPasteRepository(db, s.log)
}

func (s *Service) Create(ctx context.Context, paste *paste.Paste) (*paste.Paste, error) {
	repo := s.repo(s.db)

	newPaste, err := repo.Create(ctx, paste)

	if err != nil {
		return nil, err
	}

	return newPaste, nil
}

func (s *Service) GetByID(ctx context.Context, pasteUuid string, userId int64) (*paste.Paste, error) {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)

	if err != nil {
		return nil, err
	}

	res, err := repo.GetByID(ctx, parsedUuid, userId)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetPastes(ctx context.Context) ([]*paste.Paste, error) {
	repo := s.repo(s.db)

	userId := ctx.Value("user_id").(int64)
	pastes, err := repo.GetPastes(ctx, userId)

	if err != nil {
		return nil, err
	}

	return pastes, nil
}

func (s *Service) Delete(ctx context.Context, pasteUuid string) error {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(pasteUuid)
	userId := ctx.Value("user_id").(int64)

	err = repo.Delete(ctx, parsedUuid, userId)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Update(ctx context.Context, in UpdatePasteInput, pasteUuid string) (*paste.Paste, error) {
	parsedUuid, err := uuid.Parse(pasteUuid)
	userId := ctx.Value("user_id").(int64)

	if err != nil {
		return nil, err
	}

	foundPaste, err := s.GetByID(ctx, pasteUuid, userId)

	if err != nil {
		return nil, err
	}

	if foundPaste == nil {
		return nil, paste.ErrPasteNotFound
	}

	err = s.repo(s.db).Update(ctx, &paste.Paste{
		Uuid:       parsedUuid,
		Title:      in.Title,
		Syntax:     paste.Syntax(in.Syntax),
		Visibility: paste.Visibility(in.Visibility),
		MaxViews:   in.MaxViews,
		ExpiresAt:  in.ExpireAt,
	})

	return nil, nil
}
