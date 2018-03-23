package schema

import (
	"context"

	"github.com/go-ggz/ggz/helper"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/loader"

	"gopkg.in/nicksrandall/dataloader.v5"
)

func getUserFromLoader(ctx context.Context, id interface{}) (*model.User, error) {
	key := helper.GetCacheKey("user", id)
	userCache, err := loader.UserIDCache.Load(ctx, dataloader.StringKey(key))()

	if err != nil {
		return nil, err
	}

	return userCache.(*model.User), nil
}
