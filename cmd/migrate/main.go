package main

import (
	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/db"
)

func main() {
	cfg := config.GetConfig()
	db.RunMigrate(&cfg.DB)
}
