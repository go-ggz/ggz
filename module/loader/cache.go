package loader

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/loader/lru"
	"github.com/go-ggz/ggz/module/loader/memory"

	"gopkg.in/nicksrandall/dataloader.v5"
)

var (
	// Cache for dataloader
	Cache dataloader.Cache
	// UserIDCache for user cache from ID
	UserCache *dataloader.Loader
)

const sep = ":"

func initLoader() {
	UserCache = dataloader.NewBatchedLoader(userBatch, dataloader.WithCache(Cache))
}

func getCacheKey(module string, id interface{}) string {
	var str string
	switch v := id.(type) {
	case int64:
		str = strconv.FormatInt(v, 10)
	case string:
		str = v
	}
	return module + sep + str
}

func getCacheID(key string) (interface{}, error) {
	strs := strings.Split(key, sep)

	return strs[1], nil
}

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

func userBatch(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var results []*dataloader.Result
	id, _ := getCacheID(keys[0].String())

	user, err := model.GetUserByID(id.(int64))

	results = append(results, &dataloader.Result{
		Data:  user,
		Error: err,
	})

	return results
}

// GetUserFromLoader get user cache
func GetUserFromLoader(ctx context.Context, id interface{}) (*model.User, error) {
	key := getCacheKey("user", id)
	userCache, err := UserCache.Load(ctx, dataloader.StringKey(key))()

	if err != nil {
		return nil, err
	}

	return userCache.(*model.User), nil
}
