package graphql

import (
	"github.com/go-ggz/ggz/schema"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

// Handler initializes the prometheus middleware.
func Handler() gin.HandlerFunc {
	// Creates a GraphQL-go HTTP handler with the defined schema
	h := handler.New(&handler.Config{
		Schema: &schema.Schema,
		Pretty: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
