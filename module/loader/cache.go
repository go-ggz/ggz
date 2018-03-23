package loader

import (
	"context"

	"github.com/go-ggz/ggz/helper"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/loader/lru"
	"github.com/go-ggz/ggz/module/loader/memory"

	"gopkg.in/nicksrandall/dataloader.v5"
)

var (
	// Cache for dataloader
	Cache dataloader.Cache
	// UserIDCache for user cache from ID
	UserIDCache *dataloader.Loader
)

// NewEngine for initialize cache engine
func NewEngine(driver, prefix string, expire int) error {
	switch driver {
	case "lru":
		Cache = lru.NewEngine(prefix)
	case "memory":
		Cache = memory.NewEngine(prefix, expire)
	default:
		Cache = dataloader.NewCache()
	}

	// load cache
	initLoader()

	return nil
}

func initLoader() {
	UserIDCache = dataloader.NewBatchedLoader(userBatch, dataloader.WithCache(Cache))
}

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
