package main

import (
	"os"
	"time"

	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/version"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

func authorList() []*cli.Author {
	return []*cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
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
	}
}

func globalCommands() []*cli.Command {
	return []*cli.Command{
		Server(),
		Health(),
		Mail(),
	}
}

func globalBefore() cli.BeforeFunc {
	return func(c *cli.Context) error {
		setupLogger()
		return nil
	}
}

func main() {
	if env := os.Getenv("GGZ_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:      "gzz server",
		Usage:     "shorten url service",
		Copyright: "Copyright (c) 2019 Bo-Yi Wu",
		Version:   version.PrintCLIVersion(),
		Compiled:  time.Now(),
		Authors:   authorList(),
		Flags:     globalFlags(),
		Commands:  globalCommands(),
		Before:    globalBefore(),
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
		log.Fatal().Err(err).Msg("can't run app")
	}
}
