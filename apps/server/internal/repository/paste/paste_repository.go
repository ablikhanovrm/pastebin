package paste

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
)

type PasteRepository struct {
	q *dbgen.Queries
}

func NewPasteRepository(q *dbgen.Queries) *PasteRepository {
	return &PasteRepository{q: q}
}
