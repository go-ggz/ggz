package schema

import (
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/loader"
	"github.com/go-ggz/ggz/pkg/errors"
	"github.com/go-ggz/ggz/pkg/helper"

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

var queryMe = graphql.Field{
	Name:        "QueryMe",
	Description: "Query Cureent User",
	Type:        userType,
	Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
		user := helper.GetUserDataFromModel(p.Context)
		if user == nil {
			return nil, errors.EUnauthorized(errorYouAreNotLogin, nil)
		}

		return loader.GetUserFromLoader(p.Context, user.ID)
	},
}

func init() {
	userType.AddFieldConfig("urls", &graphql.Field{
		Type: graphql.NewList(shortenType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			o, ok := p.Source.(*model.User)

			if !ok {
				return nil, errors.ENotFound(errorUserNotFound, nil)
			}

			return model.GetShortenURLs(o.ID, 0, 10, "")
		},
	})
}
