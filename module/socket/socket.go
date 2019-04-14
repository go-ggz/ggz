package socket

import (
	// "context"
	// "errors"
	// "fmt"
	// "net/http"

	// "github.com/go-ggz/ggz/config"
	// "github.com/go-ggz/ggz/helper"
	// "github.com/go-ggz/ggz/router/middleware/auth/auth0"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	// "github.com/rs/zerolog/log"
)

// Server for socket server
var Server *socketio.Server
var err error
var key = "user"

// Test for testing websocket
type Test struct {
	A int    `json:"abc"`
	B string `json:"def"`
}

// NewEngine for socket server
func NewEngine() error {
	// Server, err = socketio.NewServer(nil)
	// if err != nil {
	// 	log.Error().Err(err).Msg("can't create socker server.")
	// 	return err
	// }

	// Server.SetAllowRequest(func(r *http.Request) error {
	// 	token := r.URL.Query().Get("token")

	// 	if token == "" {
	// 		return errors.New("Required authorization token not found")
	// 	}

	// 	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
	// 		return auth0.ParseRSAPublicKeyFromPEM()
	// 	})

	// 	if err != nil {
	// 		return fmt.Errorf("Error parsing token: %v", err)
	// 	}

	// 	if jwt.SigningMethodHS256.Alg() != parsedToken.Header["alg"] {
	// 		message := fmt.Sprintf("Expected %s signing method but token specified %s",
	// 			jwt.SigningMethodHS256.Alg(),
	// 			parsedToken.Header["alg"])
	// 		return fmt.Errorf("Error validating token algorithm: %s", message)
	// 	}

	// 	if !parsedToken.Valid {
	// 		return errors.New("Token is invalid")
	// 	}

	// 	// If we get here, everything worked and we can set the
	// 	// user property in context.
	// 	newRequest := r.WithContext(context.WithValue(r.Context(), config.ContextKeyUser, parsedToken))
	// 	// Update the current request with the new context information.
	// 	*r = *newRequest

	// 	return nil
	// })

	// Server.On("connection", func(so socketio.Socket) {
	// 	user := helper.GetUserDataFromToken(so.Request().Context())
	// 	room := user["email"].(string)
	// 	so.Join(room)

	// 	so.On("chat message", func(msg string) {
	// 		so.BroadcastTo(room, "chat message", Test{
	// 			A: 1,
	// 			B: "100",
	// 		})
	// 	})

	// 	so.On("chat message with ack", func(msg string) string {
	// 		return msg
	// 	})

	// 	so.On("disconnection", func() {
	// 	})
	// })

	// Server.On("error", func(so socketio.Socket, err error) {
	// 	log.Error().Err(err).Msg("socker server error.")
	// })

	return nil
}

// Handler initializes the prometheus middleware.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", origin)
		Server.ServeHTTP(c.Writer, c.Request)
	}
}
