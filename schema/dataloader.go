package schema

import (
	"context"

	"github.com/go-ggz/ggz/helper"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/loader"

	"gopkg.in/nicksrandall/dataloader.v5"
)

func userBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	id, _ := helper.GetCacheID(keys[0].String())

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
