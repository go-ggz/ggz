package schema

import (
	"github.com/graphql-go/graphql"
)

var rootQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "RootQuery",
		Description: "Root Query",
		Fields: graphql.Fields{
			"queryURLMetadata":   &queryURLMeta,
			"queryShortenURL":    &queryShortenURL,
			"queryAllShortenURL": &queryAllShortenURL,
		},
	})

var rootMutation = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "RootMutation",
		Description: "Root Mutation",
		Fields: graphql.Fields{
			"createShortenURL": &createShortenURL,
		},
	})

// Schema is the GraphQL schema served by the server.
var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
