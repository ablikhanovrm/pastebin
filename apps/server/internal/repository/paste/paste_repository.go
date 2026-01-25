package paste

import (
	"context"

	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
)

type PasteRepository interface {
	FindByID(ctx context.Context, id int64) (*dbgen.Paste, error)
	Create(ctx context.Context, u *dbgen.Paste) (*dbgen.Paste, error)
}

type SqlcPasteRepository struct {
	db dbgen.DBTX
}

func NewSqlcPasteRepository(db dbgen.DBTX) *SqlcPasteRepository {
	return &SqlcPasteRepository{db: db}
}

func (*SqlcPasteRepository) FindByID(ctx context.Context, id int64) (*dbgen.Paste, error) {
	return nil, nil
}

func (*SqlcPasteRepository) Create(ctx context.Context, u *dbgen.Paste) (*dbgen.Paste, error) {
	return nil, nil
}
