package migrations

import (
	"embed"

	_ "github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed ../migrations/*.sql
var FS embed.FS
