package repository

import (
	"context"
	"fmt"

	"github.com/ablikhanovrm/pastebin/internal/config"
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	Pool    *pgxpool.Pool
	Queries *dbgen.Queries
}

func NewPostgresStorage(cfg *config.DatabaseConfig) (*PostgresStorage, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	queries := dbgen.New(pool)

	return &PostgresStorage{
		Pool:    pool,
		Queries: queries,
	}, nil
}
