package repository

import (
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
)

type Repository struct {
	q *dbgen.Queries
}

func NewRepository(q *dbgen.Queries) *Repository {
	return &Repository{q}
}
