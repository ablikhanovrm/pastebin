package logging

import (
	"os"

	"github.com/rs/zerolog"
)

func New(appName string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	level := zerolog.InfoLevel

	if l, err := zerolog.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		level = l
	}

	zerolog.SetGlobalLevel(level)

	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("app", appName).
		Logger()
}
