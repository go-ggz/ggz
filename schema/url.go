package schema

import (
	"github.com/go-ggz/ggz/module/meta"

	"github.com/graphql-go/graphql"
)

var urlType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "URL",
	Description: "URL Type",
	Fields: graphql.Fields{
		"Scheme": &graphql.Field{
			Type: graphql.String,
		},
		"Opaque": &graphql.Field{
			Type: graphql.String,
		},
		"User": &graphql.Field{
			Type: graphql.String,
		},
		"Host": &graphql.Field{
			Type: graphql.String,
		},
		"Path": &graphql.Field{
			Type: graphql.String,
		},
		"RawPath": &graphql.Field{
			Type: graphql.String,
		},
		"ForceQuery": &graphql.Field{
			Type: graphql.Boolean,
		},
		"RawQuery": &graphql.Field{
			Type: graphql.String,
		},
		"Fragment": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var urlMetaType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "URLMeta",
	Description: "URL Meta Type",
	Fields: graphql.Fields{
		"title": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"image": &graphql.Field{
			Type: graphql.String,
		},
		"time": &graphql.Field{
			Type: graphql.DateTime,
		},
		"video_width": &graphql.Field{
			Type: graphql.Int,
		},
		"video_height": &graphql.Field{
			Type: graphql.Int,
		},
		"url": &graphql.Field{
			Type: urlType,
		},
	},
})

var queryURLMeta = graphql.Field{
	Name:        "QueryURLMeta",
	Description: "Query URL Metadata",
	Type:        urlMetaType,
	Args: graphql.FieldConfigArgument{
		"url": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
		url, _ := p.Args["url"].(string)

		return meta.FetchData(url)
	},
}
