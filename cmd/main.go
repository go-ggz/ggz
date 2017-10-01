package main

import (
	"os"
	"time"

	"github.com/go-ggz/ggz/config"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"
)

// Version set at compile-time
var Version = "v1.0.0-dev"

func main() {
	if env := os.Getenv("GGZ_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:      "gzz",
		Usage:     "shorten url service",
		Copyright: "Copyright (c) 2017 Bo-Yi Wu",
		Version:   Version,
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
				Value:       false,
				Usage:       "Activate debug information",
				EnvVars:     []string{"GGZ_DEBUG"},
				Destination: &config.Debug,
				Hidden:      true,
			},
		},

		Before: func(c *cli.Context) error {
			logrus.SetOutput(os.Stdout)

			if config.Debug {
				logrus.SetLevel(logrus.DebugLevel)
			} else {
				logrus.SetLevel(logrus.InfoLevel)
			}

			return nil
		},

		Commands: []*cli.Command{
			Server(),
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
