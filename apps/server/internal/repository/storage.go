package repository

import (
	"context"
	"fmt"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/jackc/pgx/v5"
)

func NewPostgresStorage(cfg *config.DatabaseConfig) (*pgx.Conn, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		return nil, err
		// fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		// os.Exit(1)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
