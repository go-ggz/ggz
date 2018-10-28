// +build windows

package main

import (
	"net/http"
)

func startServer(s *http.Server) error {
	if s.TLSConfig == nil {
		return s.ListenAndServe()
	} else {
		return s.ListenAndServeTLS("", "")
	}
}
