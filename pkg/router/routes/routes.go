package routes

import (
	"net/http"
	"path"
	"regexp"

	"github.com/go-ggz/ggz/api"
	"github.com/go-ggz/ggz/assets"
	"github.com/go-ggz/ggz/pkg/config"
	"github.com/go-ggz/ggz/pkg/middleware/auth"
	"github.com/go-ggz/ggz/pkg/middleware/graphql"
	"github.com/go-ggz/ggz/pkg/middleware/header"
	"github.com/go-ggz/ggz/pkg/model"
	"github.com/go-ggz/ggz/pkg/module/loader"
	"github.com/go-ggz/ggz/pkg/module/metrics"
	"github.com/go-ggz/ggz/pkg/module/storage"
	"github.com/go-ggz/ggz/pkg/router"

	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	rxURL = regexp.MustCompile(`^/(socket.io|graphql).*`)
)

// GlobalInit is for global configuration reload-able.
func GlobalInit() {
	if err := model.NewEngine(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize ORM engine.")
	}

	// initial socket module
	// if err := socket.NewEngine(); err != nil {
	// 	log.Fatal().Err(err).Msg("Failed to initialize Socket IO engine")
	// }

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
	if config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	c := metrics.NewCollector()
	prometheus.MustRegister(c)

	e := gin.New()

	e.Use(gin.Recovery())
	e.Use(logger.SetLogger(logger.Config{
		UTC:            true,
		SkipPathRegexp: rxURL,
	}))
	// e.Use(gzip.Gzip(gzip.DefaultCompression))
	e.Use(header.Options)
	e.Use(header.Secure)
	e.Use(middleware...)

	if config.Server.Pprof {
		pprof.Register(
			e,
			path.Join(config.Server.Root, "debug", "pprof"),
		)
	}

	// redirect to vue page
	e.NoRoute(gzip.Gzip(gzip.DefaultCompression), api.Index)

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
			"/public",
			assets.Load(),
		)

		root.GET("", gzip.Gzip(gzip.DefaultCompression), api.Index)
		root.GET("/favicon.ico", api.Favicon)
		root.GET("/metrics", router.Metrics(config.Prometheus.AuthToken))
		root.GET("/healthz", api.Heartbeat)
		root.GET("/assets/*name", gzip.Gzip(gzip.DefaultCompression), assets.ViewHandler())

		v := e.Group("/v1")
		v.Use(auth.Check())
		{
			v.POST("/url/meta", api.URLMeta)
			v.POST("/s", api.CreateShortenURL)
		}

		g := e.Group("/graphql")
		g.Use(auth.Check())
		{
			g.POST("", graphql.Handler())
			if config.Server.GraphiQL {
				g.GET("", graphql.Handler())
			}
		}

		// socket connection
		// root.GET("/socket.io/", socket.Handler())
		// root.POST("/socket.io/", socket.Handler())
		// root.Handle("WS", "/socket.io", socket.Handler())
		// root.Handle("WSS", "/socket.io", socket.Handler())
	}

	return e
}

// LoadRedirct initializes the routing of the shorten URL application.
func LoadRedirct(middleware ...gin.HandlerFunc) http.Handler {
	if config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.Use(gin.Recovery())
	e.Use(logger.SetLogger(logger.Config{
		UTC:            true,
		SkipPathRegexp: rxURL,
	}))
	e.Use(header.Options)
	e.Use(header.Secure)
	e.Use(middleware...)

	if config.Server.Pprof {
		pprof.Register(
			e,
			path.Join(config.Server.Root, "debug", "pprof"),
		)
	}

	// 404 not found
	e.NoRoute(api.NotFound)

	// default route /
	root := e.Group(config.Server.Root)
	{
		root.GET("", api.Index)
		root.GET("/:slug", api.RedirectURL)
	}

	return e
}
