package schema

import (
	"context"

	"github.com/go-ggz/ggz/helper"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/loader"

	"gopkg.in/nicksrandall/dataloader.v4"
)

func userBatch(ctx context.Context, keys []interface{}) []*dataloader.Result {
	var results []*dataloader.Result
	id, _ := helper.GetCacheID(keys[0].(string))

	user, err := model.GetUserByID(id.(int64))

	results = append(results, &dataloader.Result{
		Data:  user,
		Error: err,
	})

	return results
}

var (
	userLoader = dataloader.NewBatchedLoader(userBatch, dataloader.WithCache(loader.Cache))
)
