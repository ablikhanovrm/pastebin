package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrate(cfg *config.DatabaseConfig) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"pgx",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[len(os.Args)-1] {
	case "up":
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
