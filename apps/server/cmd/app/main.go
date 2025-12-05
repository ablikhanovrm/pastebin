package main

import (
	"github.com/ablikhanovrm/pastebin/internal/app"
)

const configPath = "../../configs/main.yaml"

func main() {
	app.Run(configPath)
}
