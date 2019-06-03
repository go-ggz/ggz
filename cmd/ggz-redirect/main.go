package main

import (
	"os"
	"strings"
	"time"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/version"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

func setupLogging() {
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

func main() {
	if env := os.Getenv("GGZ_ENV_FILE"); env != "" {
		if err := godotenv.Load(env); err != nil {
			log.Fatal().Err(err).Msg("Cannot start load config from env")
		}
	}

	app := &cli.App{
		Name:      "gzz redirect",
		Usage:     "redirect service",
		Copyright: "Copyright (c) 2018 Bo-Yi Wu",
		Version:   version.PrintCLIVersion(),
		Compiled:  time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Bo-Yi Wu",
				Email: "appleboy.tw@gmail.com",
			},
		},

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "debug",
				Value:       true,
				Usage:       "Activate debug information",
				EnvVars:     []string{"GGZ_SERVER_DEBUG"},
				Destination: &config.Server.Debug,
			},
			&cli.BoolFlag{
				Name:        "log-color",
				Value:       true,
				Usage:       "enable colored logging",
				EnvVars:     []string{"GGZ_LOGS_COLOR"},
				Destination: &config.Logs.Color,
			},
			&cli.BoolFlag{
				Name:        "log-pretty",
				Value:       true,
				Usage:       "enable pretty logging",
				EnvVars:     []string{"GGZ_LOGS_PRETTY"},
				Destination: &config.Logs.Pretty,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Value:       "info",
				Usage:       "set logging level",
				EnvVars:     []string{"GGZ_LOGS_LEVEL"},
				Destination: &config.Logs.Level,
			},
		},

		Before: func(c *cli.Context) error {
			setupLogging()

			return nil
		},

		Commands: []*cli.Command{
			Server(),
			Health(),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
