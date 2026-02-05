package paste

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	authrepo "github.com/ablikhanovrm/pastebin/internal/repository/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Service struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func (s *Service) repo(db dbgen.DBTX) *authrepo.SqlcAuthRepository {
	return authrepo.NewSqlcAuthRepository(db, s.log)
}

func NewPasteService(db *pgxpool.Pool, log zerolog.Logger) *Service {
	return &Service{db: db, log: log}
}
