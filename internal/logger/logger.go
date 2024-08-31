package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger(logLevel string) {
	l, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Fatal().Msgf("log level=%q is not valid", logLevel)
	}
	zerolog.SetGlobalLevel(l)
	log.Logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
}
