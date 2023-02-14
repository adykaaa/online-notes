package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var logger zerolog.Logger

func New(level string) zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	loglevel, err := zerolog.ParseLevel(level)
	if err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.SetGlobalLevel(loglevel)

	logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		CallerWithSkipFrameCount(2).
		Logger()

	zerolog.DefaultContextLogger = &logger

	return logger
}
