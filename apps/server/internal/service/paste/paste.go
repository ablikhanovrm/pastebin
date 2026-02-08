package paste

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/ablikhanovrm/pastebin/internal/models/paste"
	authrepo "github.com/ablikhanovrm/pastebin/internal/repository/auth"
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
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewPasteService(db *pgxpool.Pool, log zerolog.Logger) *Service {
	return &Service{db: db, log: log}
}

// repo helper
func (s *Service) repo(db dbgen.DBTX) *authrepo.SqlcAuthRepository {
	return authrepo.NewSqlcAuthRepository(db, s.log)
}

func (s *Service) Create(ctx context.Context, paste *paste.Paste) (*paste.Paste, error) {
	return nil, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*paste.Paste, error) {
	return nil, nil
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
