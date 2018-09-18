package main

import (
	"os"
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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
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
		Name:      "gzz",
		Usage:     "shorten url service",
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
				EnvVars:     []string{"GGZ_DEBUG"},
				Destination: &config.Debug,
			},
			&cli.BoolFlag{
				Name:        "color",
				Value:       true,
				Usage:       "Enable pprof debugging server",
				EnvVars:     []string{"GGZ_LOGS_COLOR"},
				Destination: &config.Logs.Color,
			},
			&cli.BoolFlag{
				Name:        "pretty",
				Value:       true,
				Usage:       "Enable pprof debugging server",
				EnvVars:     []string{"GGZ_LOGS_PRETTY"},
				Destination: &config.Logs.Pretty,
			},
		},

		Before: func(c *cli.Context) error {
			setupLogging()

			return nil
		},

		Commands: []*cli.Command{
			Server(),
			Ping(),
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
