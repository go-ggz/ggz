package prometheus

import (
	"errors"
	"fmt"

	"github.com/go-ggz/ggz/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// errInvalidToken is returned when the api request token is invalid.
	errInvalidToken = errors.New("Invalid or missing token")
)

// Handler initializes the prometheus middleware.
func Handler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		token := config.Prometheus.AuthToken

		if token == "" {
			h.ServeHTTP(c.Writer, c.Request)
			return
		}

		header := c.Request.Header.Get("Authorization")

		if header == "" {
			c.String(401, errInvalidToken.Error())
			return
		}

		bearer := fmt.Sprintf("Bearer %s", token)

		if header != bearer {
			c.String(401, errInvalidToken.Error())
			return
		}

		h.ServeHTTP(c.Writer, c.Request)
	}
}
