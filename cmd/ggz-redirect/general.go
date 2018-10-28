// +build !windows

package main

import (
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
)

func startServer(s *http.Server) error {
	return gracehttp.Serve(s)
}
