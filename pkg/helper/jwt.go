package helper

import (
	"context"

	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/model"

	"github.com/dgrijalva/jwt-go"
)

// GetUserDataFromToken from jwt parse token
func GetUserDataFromToken(ctx context.Context) jwt.MapClaims {
	if _, ok := ctx.Value("user").(*jwt.Token); !ok {
		return jwt.MapClaims{}
	}

	return ctx.Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
}

// GetUserDataFromModel from user model
func GetUserDataFromModel(ctx context.Context) *model.User {
	if _, ok := ctx.Value(config.ContextKeyUser).(*model.User); !ok {
		return nil
	}

	return ctx.Value(config.ContextKeyUser).(*model.User)
}
