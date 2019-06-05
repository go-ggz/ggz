package schema

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/pkg/config"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
)

// T for graphql testing schema
type T struct {
	Query    string
	Schema   graphql.Schema
	Expected *graphql.Result
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
	if len(result.Errors) != len(test.Expected.Errors) {
		t.Fatalf("Unexpected errors, Diff: %v", testutil.Diff(test.Expected.Errors, result.Errors))
	}

	if len(test.Expected.Errors) > 0 &&
		result.Errors[0].Message != test.Expected.Errors[0].Message {
		t.Fatalf("Unexpected error message, Diff: %v", testutil.Diff(test.Expected.Errors, result.Errors))
	}
}

func newContextWithUser(ctx context.Context, u *model.User) context.Context {
	return context.WithValue(ctx, config.ContextKeyUser, u)
}
