package logging

import (
	"os"

	"github.com/rs/zerolog"
)

func New(appName string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("app", appName).
		Logger()
}
