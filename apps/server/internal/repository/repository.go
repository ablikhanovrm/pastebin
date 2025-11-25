package repository

import "github.com/jackc/pgx/v5"

type Repository struct {
}

func (r *Repository) NewRepository(*pgx.Conn) *Repository {
	return &Repository{}
}
