package auth0

import (
	"errors"
	"net/http"

	"github.com/go-ggz/ggz/assets"
	"github.com/go-ggz/ggz/config"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func errorHandler(w http.ResponseWriter, r *http.Request, err string) {
}

// Check initializes the auth0 middleware.
func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				var reader []byte
				var err error

				if config.Auth0.Key != "" {
					reader = []byte(config.Auth0.Key)
				} else {
					reader, err = assets.ReadFile(config.Auth0.PemPath)
					if err != nil {
						logrus.Warnf("Failed to read builtin %s template. %s", reader, err)
						return nil, errors.New("Failed to read builtin auth0 pem file")
					}
				}

				return jwt.ParseRSAPublicKeyFromPEM(reader)
			},
			SigningMethod: jwt.SigningMethodRS256,
			ErrorHandler:  errorHandler,
			Debug:         config.Auth0.Debug,
		})

		err := jwtMiddleware.CheckJWT(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			logrus.Errorf("JWT Error: %s", err.Error())
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "You don't have permission",
				},
			)
			return
		}

		user := c.Request.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
		c.Set("email", user["email"])
		c.Set("email_verified", user["email_verified"])
		c.Set("name", user["name"])
	}
}
