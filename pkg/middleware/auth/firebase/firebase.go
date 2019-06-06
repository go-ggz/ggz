package firebase

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-ggz/ggz/assets"
	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/model"

	"firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

var (
	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
)

func getFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	if authHeader == "" {
		return "", ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}

// Check initializes the firebase middleware.
func Check() gin.HandlerFunc {
	credentials, err := assets.ReadSource("/firebase/serviceAccountKey.json")
	if err != nil {
		log.Fatal().Err(err).Msg("can't load credentials")
	}
	ctx := context.Background()
	opt := option.WithCredentialsJSON(credentials)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initial firebase app")
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initial firebase client")
	}

	return func(c *gin.Context) {
		token, err := getFromHeader(c, "Authorization")

		if err != nil {
			c.Next()
		} else {
			userData, err := client.VerifyIDToken(ctx, token)
			if err != nil {
				log.Error().Err(err).Msg("verify firebase token error.")
				c.AbortWithStatusJSON(
					http.StatusOK,
					gin.H{
						"data": nil,
						"errors": []map[string]interface{}{
							{
								"message": "token expire or parse error",
							},
						},
					},
				)
				return
			}

			// check user exist
			user, err := model.GetUserByEmail(userData.Claims["email"].(string))

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
					Email:    userData.Claims["email"].(string),
					IsActive: userData.Claims["email_verified"].(bool),
				}

				if v, ok := userData.Claims["name"]; ok {
					user.FullName = v.(string)
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
