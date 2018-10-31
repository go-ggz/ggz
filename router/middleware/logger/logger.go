package logger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// SetLogger initializes the logging middleware.
func SetLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		msg := "Request"
		if c.Errors.String() != "" {
			msg = c.Errors.String()
		}

		logger := log.With().
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			Str("user-agent", c.Request.UserAgent()).
			Logger()

		switch {
		case c.Writer.Status() >= 400 && c.Writer.Status() < 500:
			{
				logger.Warn().
					Msg(msg)
			}
		case c.Writer.Status() >= 500:
			{
				logger.Error().
					Msg(msg)
			}
		default:
			logger.Info().
				Msg(msg)
		}
	}
}
