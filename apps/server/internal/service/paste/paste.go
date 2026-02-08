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
	Create(ctx context.Context, paste *paste.Paste) (*paste.Paste, error)
	GetByID(ctx context.Context, id int64) (*paste.Paste, error)
	GetAll(ctx context.Context) ([]*paste.Paste, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, paste *paste.Paste) (*paste.Paste, error)
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
	return nil, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*paste.Paste, error) {
	repo := s.repo(s.db)

	parsedUuid, err := uuid.Parse(id)

	if err != nil {
		return nil, err
	}

	res, err := repo.GetByID(ctx, parsedUuid)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*paste.Paste, error) {
	return nil, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return nil
}

func (s *Service) Update(ctx context.Context, paste *paste.Paste) (*paste.Paste, error) {
	return nil, nil
}
