package auth

import (
	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/middleware/auth/auth0"
	"github.com/go-ggz/ggz/pkg/middleware/auth/firebase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Check initializes the auth0 middleware.
func Check() gin.HandlerFunc {
	switch config.Auth.Driver {
	case "auth0":
		return auth0.Check()
	case "firebase":
		return firebase.Check()
	default:
		log.Fatal().Msgf("Can't find the auth driver: %s", config.Auth.Driver)
	}

	return nil
}
