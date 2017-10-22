package helper

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

// GetUserData from jwt parse token
func GetUserData(ctx context.Context) jwt.MapClaims {
	return ctx.Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
}
