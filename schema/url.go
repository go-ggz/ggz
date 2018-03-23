package schema

import (
	"fmt"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/helper"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/meta"
	"github.com/go-ggz/ggz/web"

	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

var shortenType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ShortenType",
	Description: "Shorten URL Type",
	Fields: graphql.Fields{
		"slug": &graphql.Field{
			Type: graphql.String,
		},
		"user": &graphql.Field{
			Type: userType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				source := p.Source

				if o, ok := source.(*model.Shorten); ok {
					if o.User != nil {
						return o.User, nil
					}

					return getUserFromLoader(p.Context, o.UserID)
				}

				return nil, fmt.Errorf("source is empty")
			},
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.DateTime,
		},
		"hits": &graphql.Field{
			Type: graphql.Int,
		},
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
	},
})

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

var createShortenURL = graphql.Field{
	Name:        "CreateShortenURL",
	Description: "Create Shorten URL",
	Type:        shortenType,
	Args: graphql.FieldConfigArgument{
		"url": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
		url, _ := p.Args["url"].(string)
		user := helper.GetUserDataFromModel(p.Context)

		row, err := model.GetShortenFromURL(url)

		if model.IsErrURLExist(err) {
			return row, nil
		}

		if err != nil {
			return nil, err
		}

		row, err = model.NewShortenURL(url, config.Server.ShortenSize, user)

		if err != nil {
			return nil, err
		}

		// upload QRCode image.
		go func(slug string) {
			if err := web.QRCodeGenerator(slug); err != nil {
				logrus.Errorln(err)
			}
		}(row.Slug)

		return row, nil
	},
}

var queryShortenURL = graphql.Field{
	Name:        "QueryShortenURL",
	Description: "Query Shorten URL",
	Type:        shortenType,
	Args: graphql.FieldConfigArgument{
		"slug": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
		slug, _ := p.Args["slug"].(string)

		return model.GetShortenBySlug(slug)
	},
}

var queryAllShortenURL = graphql.Field{
	Name:        "QueryAllShortenURL",
	Description: "Query All Shorten URL",
	Type:        graphql.NewList(shortenType),
	Args: graphql.FieldConfigArgument{
		"userID": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"page": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 1,
		},
		"pageSize": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 10,
		},
	},
	Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
		id, _ := p.Args["userID"].(int)
		page, _ := p.Args["page"].(int)
		pageSize, _ := p.Args["pageSize"].(int)
		userID := int64(id)

		return model.GetShortenURLs(userID, page, pageSize, "")
	},
}
