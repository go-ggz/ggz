package helper

import (
	"context"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/model"

	"github.com/dgrijalva/jwt-go"
)

// GetUserDataFromToken from jwt parse token
func GetUserDataFromToken(ctx context.Context) jwt.MapClaims {
	return ctx.Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
}

// GetUserDataFromModel from user model
func GetUserDataFromModel(ctx context.Context) *model.User {
	return ctx.Value(config.ContextKeyUser).(*model.User)
}
