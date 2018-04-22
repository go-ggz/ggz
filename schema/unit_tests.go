package schema

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/model"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

// T for graphql testing schema
type T struct {
	Query    string
	Schema   graphql.Schema
	Expected interface{}
}

func testGraphql(test T, p graphql.Params, t *testing.T) {
	result := graphql.Do(p)
	if len(result.Errors) > 0 {
		t.Fatalf("wrong result, unexpected errors: %v", result.Errors)
	}
	if !reflect.DeepEqual(result, test.Expected) {
		t.Fatalf("wrong result, query: %v, graphql result diff: %v", test.Query, testutil.Diff(test.Expected, result))
	}
}

func testGraphqlErr(test T, p graphql.Params, t *testing.T) {
	result := graphql.Do(p)
	if len(result.Errors) == 0 {
		t.Fatalf("missing errors, expected errors: %v", result.Errors)
	}
}

func newContextWithUser(ctx context.Context, u *model.User) context.Context {
	return context.WithValue(ctx, config.ContextKeyUser, u)
}
