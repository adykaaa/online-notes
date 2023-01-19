package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var logger zerolog.Logger

func NewLogger(level string) zerolog.Logger {

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	//if LOG_LEVEL env var cannot be parsed, default to INFO level
	loglevel, err := zerolog.ParseLevel(level)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	zerolog.SetGlobalLevel(loglevel)

	logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	return logger
}
