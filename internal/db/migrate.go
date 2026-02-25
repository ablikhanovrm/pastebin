package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrate(cfg *config.DatabaseConfig) {
	fmt.Println("RUN MIGRATE", cfg)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, cfg.SslMode)

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	abs, _ := filepath.Abs("./migrations")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+abs,
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
