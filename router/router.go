package router

import (
	"net/http"
	"path"

	"github.com/go-ggz/ggz/assets"
	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/model"
	"github.com/go-ggz/ggz/module/loader"
	"github.com/go-ggz/ggz/module/socket"
	"github.com/go-ggz/ggz/module/storage"
	"github.com/go-ggz/ggz/router/middleware/auth0"
	"github.com/go-ggz/ggz/router/middleware/graphql"
	"github.com/go-ggz/ggz/router/middleware/header"
	"github.com/go-ggz/ggz/router/middleware/prometheus"
	"github.com/go-ggz/ggz/web"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	if err := model.NewEngine(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize ORM engine.")
	}

	// initial socket module
	if err := socket.NewEngine(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Socket IO engine")
	}

	if config.QRCode.Enable {
		var err error
		storage.S3, err = storage.NewEngine()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create s3 interface")
		}

		if err := storage.S3.CreateBucket(config.Minio.Bucket, config.Minio.Region); err != nil {
			log.Fatal().Err(err).Msg("Failed to create s3 bucket")
		}
	}

	// initial dataloader cache
	if err := loader.NewEngine(config.Cache.Driver, config.Cache.Prefix, config.Cache.Expire); err != nil {
		log.Fatal().Err(err).Msg("Failed to initial dataloader.")
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
	// e.Use(gzip.Gzip(gzip.DefaultCompression))
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
		if config.Storage.Driver == "disk" {
			root.StaticFS(
				"/storage",
				gin.Dir(
					config.Storage.Path,
					false,
				),
			)
		}

		root.StaticFS(
			"/assets",
			assets.Load(),
		)

		root.GET("", web.Index)
		root.GET("/favicon.ico", web.Favicon)
		root.GET("/metrics", prometheus.Handler())
		root.GET("/healthz", web.Heartbeat)

		api := e.Group("/v1")
		api.Use(auth0.Check())
		{
			api.POST("/url/meta", web.URLMeta)
			api.POST("/s", web.CreateShortenURL)
		}

		g := e.Group("/graphql")
		g.Use(auth0.Check())
		{
			g.POST("", graphql.Handler())
		}

		// socket connection
		root.GET("/socket.io/", socket.Handler())
		root.POST("/socket.io/", socket.Handler())
		root.Handle("WS", "/socket.io", socket.Handler())
		root.Handle("WSS", "/socket.io", socket.Handler())
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
		root.GET("/:slug", web.RedirectURL)
	}

	return e
}
