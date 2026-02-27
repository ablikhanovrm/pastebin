package repository

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ablikhanovrm/pastebin/internal/config"
	dbgen "github.com/ablikhanovrm/pastebin/internal/db/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	Pool    *pgxpool.Pool
	Queries *dbgen.Queries
}

func NewPostgresStorage(cfg *config.DatabaseConfig) (*PostgresStorage, error) {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.DbName,
	}

	q := u.Query()
	q.Set("sslmode", cfg.SslMode)
	u.RawQuery = q.Encode()

	dsn := u.String()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
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
