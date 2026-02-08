package paste

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/rs/zerolog"
)

type PasteRepository interface {
	FindByID(ctx context.Context, id int64) (*dbgen.Paste, error)
	Create(ctx context.Context, u *dbgen.Paste) (*dbgen.Paste, error)
	GetAll(ctx context.Context) ([]*dbgen.Paste, error)
	Update(ctx context.Context, paste *dbgen.Paste) (*dbgen.Paste, error)
	Delete(ctx context.Context, id int64) error
}

type SqlcPasteRepository struct {
	db     dbgen.DBTX
	logger zerolog.Logger
}

func NewSqlcPasteRepository(db dbgen.DBTX, logger zerolog.Logger) *SqlcPasteRepository {
	return &SqlcPasteRepository{db: db, logger: logger}
}

func (*SqlcPasteRepository) FindByID(ctx context.Context, id int64) (*dbgen.Paste, error) {
	return nil, nil
}

func (*SqlcPasteRepository) Create(ctx context.Context, u *dbgen.Paste) (*dbgen.Paste, error) {
	return nil, nil
}

func (*SqlcPasteRepository) GetAll(ctx context.Context) ([]*dbgen.Paste, error) {
	return nil, nil
}

func (*SqlcPasteRepository) Update(ctx context.Context, paste *dbgen.Paste) (*dbgen.Paste, error) {
	return nil, nil
}

func (*SqlcPasteRepository) Delete(ctx context.Context, id int64) error {
	return nil
}
