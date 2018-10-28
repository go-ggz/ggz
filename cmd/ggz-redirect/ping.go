package main

import (
	"fmt"
	"net/http"

	"gopkg.in/urfave/cli.v2"
)

// Ping provides the sub-command to check server live.
func Ping() *cli.Command {
	return &cli.Command{
		Name:  "ping",
		Usage: "server healthy check",
		Action: func(c *cli.Context) error {
			resp, err := http.Get("http://localhost" + defaultHostAddr + "/healthz")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				return fmt.Errorf("server returned non-200 status code")
			}
			return nil
		},
	}
}
