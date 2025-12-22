package main

import (
	"github.com/ablikhanovrm/pastebin/internal/config"
	"github.com/ablikhanovrm/pastebin/internal/db"
)

const configPath = "../../configs/main.yaml"

func main() {
	cfg := config.GetConfig(configPath)
	db.RunMigrate(&cfg.DB)
}
