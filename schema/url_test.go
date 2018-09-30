package schema

import (
	"context"
	"testing"

	"github.com/go-ggz/ggz/model"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
)

func TestQueryURLMeta(t *testing.T) {
	t.Run("invaild url", func(t *testing.T) {
		test := T{
			Query: `
query queryURLMeta (
    $url: String!
) {
  queryURLMeta(url: $url) {
    url
  }
}
	  `,
			Schema: Schema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"queryURLMeta": map[string]interface{}{
						"url": "http://example.com",
					},
				},
			},
		}
		params := graphql.Params{
			Schema:        test.Schema,
			RequestString: test.Query,
			Context:       newContextWithUser(context.TODO(), nil),
			VariableValues: map[string]interface{}{
				"url": "example.com",
			},
		}
		testGraphqlErr(test, params, t)
	})
}

func TestQueryShortenURL(t *testing.T) {
	assert.NoError(t, model.PrepareTestDatabase())
	user := model.AssertExistsAndLoadBean(t, &model.User{ID: 1}).(*model.User)
	ctx := newContextWithUser(context.TODO(), user)

	t.Run("shorten url exist", func(t *testing.T) {
		test := T{
			Query: `
query queryShortenURL (
    $slug: String!
) {
  queryShortenURL(slug: $slug) {
    url
  }
}
	  `,
			Schema: Schema,
			Expected: &graphql.Result{
				Data: map[string]interface{}{
					"queryShortenURL": map[string]interface{}{
						"url": "http://example.com",
					},
				},
			},
		}
		params := graphql.Params{
			Schema:        test.Schema,
			RequestString: test.Query,
			Context:       ctx,
			VariableValues: map[string]interface{}{
				"slug": "abcdef",
			},
		}
		testGraphql(test, params, t)
	})

	t.Run("shorten url not exist", func(t *testing.T) {
		test := T{
			Query: `
query queryShortenURL (
    $slug: String!
) {
  queryShortenURL(slug: $slug) {
    url
  }
}
`,
			Schema: Schema,
		}
		params := graphql.Params{
			Schema:        test.Schema,
			RequestString: test.Query,
			Context:       ctx,
			VariableValues: map[string]interface{}{
				"slug": "1234567890",
			},
		}
		testGraphqlErr(test, params, t)
	})

}
