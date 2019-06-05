package auth0

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"

	"github.com/go-ggz/ggz/assets"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/helper"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ParseRSAPublicKeyFromPEM Parse PEM encoded PKCS1 or PKCS8 public key
func ParseRSAPublicKeyFromPEM() (*rsa.PublicKey, error) {
	var reader []byte
	var err error

	if config.Auth0.Key != "" {
		reader = []byte(config.Auth0.Key)
	} else {
		reader, err = assets.ReadSource(config.Auth0.PemPath)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to read builtin %s template.", reader)
			return nil, errors.New("Failed to read builtin auth0 pem file")
		}
	}

	return jwt.ParseRSAPublicKeyFromPEM(reader)
}

func errorHandler(w http.ResponseWriter, r *http.Request, err string) {
}

// Check initializes the auth0 middleware.
func Check() gin.HandlerFunc {
	var user *model.User
	return func(c *gin.Context) {
		jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return ParseRSAPublicKeyFromPEM()
			},
			SigningMethod: jwt.SigningMethodRS256,
			ErrorHandler:  errorHandler,
			Debug:         config.Auth0.Debug,
		})

		err := jwtMiddleware.CheckJWT(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			c.Next()
		} else {
			userClaim := helper.GetUserDataFromToken(c.Request.Context())
			if _, ok := userClaim["email"]; !ok {
				c.AbortWithStatusJSON(
					http.StatusOK,
					gin.H{
						"data": nil,
						"errors": []map[string]interface{}{
							{
								"message": "email not found.",
							},
						},
					},
				)
				return
			}

			// check user exist
			user, err = model.GetUserByEmail(userClaim["email"].(string))

			if err != nil {
				if !model.IsErrUserNotExist(err) {
					log.Error().Err(err).Msg("database error.")
					c.AbortWithStatusJSON(
						http.StatusBadRequest,
						gin.H{
							"data": nil,
							"errors": []map[string]interface{}{
								{
									"message": "database query error",
								},
							},
						},
					)
					return
				}

				// create new user
				user = &model.User{
					Email:    userClaim["email"].(string),
					FullName: userClaim["name"].(string),
					IsActive: userClaim["email_verified"].(bool),
				}
				err := model.CreateUser(user)

				if err != nil {
					log.Error().Err(err).Msg("database error.")
					c.AbortWithStatusJSON(
						http.StatusOK,
						gin.H{
							"data": nil,
							"errors": []map[string]interface{}{
								{
									"message": "can't create new user",
								},
							},
						},
					)
					return
				}
			}

			ctx := context.WithValue(c.Request.Context(), config.ContextKeyUser, user)
			c.Request = c.Request.WithContext(ctx)
		}
	}
}
