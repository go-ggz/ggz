package schema

import (
	"github.com/go-ggz/ggz/model"

	"github.com/graphql-go/graphql"
)

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "UserType",
	Description: "User Type",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"fullname": &graphql.Field{
			Type: graphql.String,
		},
		"location": &graphql.Field{
			Type: graphql.Int,
		},
		"website": &graphql.Field{
			Type: graphql.String,
		},
		"is_active": &graphql.Field{
			Type: graphql.Boolean,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
	},
})

func init() {
	userType.AddFieldConfig("urls", &graphql.Field{
		Type: graphql.NewList(shortenType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			o, ok := p.Source.(*model.User)

			if !ok {
				return nil, errMissingSource
			}

			return model.GetShortenURLs(o.ID, 0, 10, "")
		},
	})
}
