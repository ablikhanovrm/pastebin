package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrate(cfg *config.DatabaseConfig) {
	fmt.Println("RUN MIGRATE", cfg)
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
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		d,
		"postgres",
		driver,
	)
	fmt.Println("MIGRATE")
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[len(os.Args)-1] {
	case "up":
		v, dirty, _ := m.Version()
		fmt.Println("Current version:", v, "dirty:", dirty)
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migration up error: %v", err)
		}
		fmt.Println("Migrations applied")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migration down error: %v", err)
		}
		fmt.Println("Migrations reverted")
	}
}
