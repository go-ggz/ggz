package main

import (
	"fmt"
	"net/http"

	"github.com/go-ggz/ggz/pkg/config"

	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

func healthAction() cli.ActionFunc {
	return func(c *cli.Context) error {
		resp, err := http.Get("http://localhost" + config.Server.Addr + "/healthz")
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to request health check")
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Error().
				Int("code", resp.StatusCode).
				Msg("health seems to be in bad state")
			return fmt.Errorf("server returned non-200 status code")
		}
		return nil
	}
}

func healthFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "addr",
			Value:       defaultHostAddr,
			Usage:       "Address to bind the server",
			EnvVars:     []string{"GGZ_SERVER_ADDR"},
			Destination: &config.Server.Addr,
		},
	}
}

// Health provides the sub-command to perform a health check.
func Health() *cli.Command {
	return &cli.Command{
		Name:   "health",
		Usage:  "perform health checks",
		Flags:  healthFlags(),
		Action: healthAction(),
	}
}
