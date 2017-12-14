package schema

import (
	"context"

	"github.com/go-ggz/ggz/model"

	"gopkg.in/nicksrandall/dataloader.v4"
)

var userLoader *dataloader.Loader

func userBatch(ctx context.Context, keys []interface{}) []*dataloader.Result {
	var results []*dataloader.Result
	id := keys[0].(int64)
	user, err := model.GetUserByID(id)

	results = append(results, &dataloader.Result{
		Data:  user,
		Error: err,
	})

	return results
}

func init() {
	userLoader = dataloader.NewBatchedLoader(userBatch)
}
