package main

import (
	"os"
	"strings"

	"github.com/go-ggz/ggz/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger() {
	switch strings.ToLower(config.Logs.Level) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	if config.Logs.Pretty {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: !config.Logs.Color,
			},
		)
	}
}
