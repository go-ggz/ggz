package router

import (
	"net/http"
	"path"

	"github.com/go-ggz/ggz/assets"
	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/models"
	"github.com/go-ggz/ggz/router/middleware/header"
	"github.com/go-ggz/ggz/router/middleware/logger"
	"github.com/go-ggz/ggz/router/middleware/prometheus"
	"github.com/go-ggz/ggz/web"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	if err := models.NewEngine(); err != nil {
		logrus.Fatalf("Failed to initialize ORM engine: %v", err)
	}
}

// Load initializes the routing of the application.
func Load(middleware ...gin.HandlerFunc) http.Handler {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(gin.Recovery())
	e.Use(logger.SetLogger())
	e.Use(gzip.Gzip(gzip.DefaultCompression))
	e.Use(header.Options)
	e.Use(header.Secure)
	e.Use(middleware...)

	if config.Server.Pprof {
		pprof.Register(
			e,
			&pprof.Options{
				RoutePrefix: path.Join(config.Server.Root, "debug", "pprof"),
			},
		)
	}

	// 404 not found
	e.NoRoute(web.NotFound)

	// default route /
	root := e.Group(config.Server.Root)
	{
		root.StaticFS(
			"/storage",
			gin.Dir(
				config.Server.Storage,
				false,
			),
		)

		root.StaticFS(
			"/assets",
			assets.Load(),
		)

		root.GET("", web.Index)
		root.GET("/favicon.ico", web.Favicon)
		root.GET("/metrics", prometheus.Handler())
		root.POST("/s", web.CreateShortenURL)
		root.GET("/g/:slug", web.FetchShortenedURL)
	}

	return e
}

// LoadRedirct initializes the routing of the shorten URL application.
func LoadRedirct(middleware ...gin.HandlerFunc) http.Handler {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(gin.Recovery())
	e.Use(logger.SetLogger())
	e.Use(gzip.Gzip(gzip.DefaultCompression))
	e.Use(header.Options)
	e.Use(header.Secure)
	e.Use(middleware...)

	if config.Server.Pprof {
		pprof.Register(
			e,
			&pprof.Options{
				RoutePrefix: path.Join(config.Server.Root, "debug", "pprof"),
			},
		)
	}

	// 404 not found
	e.NoRoute(web.NotFound)

	// default route /
	root := e.Group(config.Server.Root)
	{
		root.GET("", web.Index)
		root.GET("/:slug", web.ShortenedURL)
	}

	return e
}
