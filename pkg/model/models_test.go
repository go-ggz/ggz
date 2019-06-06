package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePostgreSQLHostPort(t *testing.T) {
	tests := []struct {
		Name     string
		HostPort string
		Host     string
		Port     string
	}{
		{
			Name:     "ip and port",
			HostPort: "127.0.0.1:1234",
			Host:     "127.0.0.1",
			Port:     "1234",
		},
		{
			Name:     "ip",
			HostPort: "127.0.0.1",
			Host:     "127.0.0.1",
			Port:     "5432",
		},
		{
			Name:     "ipv6 and port",
			HostPort: "[::1]:1234",
			Host:     "[::1]",
			Port:     "1234",
		},
		{
			Name:     "ipv6",
			HostPort: "[::1]",
			Host:     "[::1]",
			Port:     "5432",
		},
		{
			Name:     "socket and port",
			HostPort: "/tmp/pg.sock:1234",
			Host:     "/tmp/pg.sock",
			Port:     "1234",
		},
		{
			Name:     "socket",
			HostPort: "/tmp/pg.sock",
			Host:     "/tmp/pg.sock",
			Port:     "5432",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()
			host, port := parsePostgreSQLHostPort(tt.HostPort)
			assert.Equal(t, tt.Host, host)
			assert.Equal(t, tt.Port, port)
		})
	}
}
