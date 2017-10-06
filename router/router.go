package router

import (
	"net/http"
	"os"
	"path"

	"github.com/go-ggz/ggz/assets"
	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/models"
	"github.com/go-ggz/ggz/modules/minio"
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

	if config.QRCode.Enable && config.Storage.Driver == "disk" {
		storage := path.Join(config.Storage.Path, config.QRCode.Bucket)
		if err := os.MkdirAll(storage, os.ModePerm); err != nil {
			logrus.Fatalf("Failed to create storage folder: %v", err)
		}
	}

	if config.QRCode.Enable && config.Storage.Driver == "s3" {
		minio.NewEngine(
			config.Minio.EndPoint,
			config.Minio.AccessID,
			config.Minio.SecretKey,
			config.Minio.SSL,
		)

		if err := minio.S3.MakeBucket(config.Minio.Bucket, config.Minio.Region); err != nil {
			logrus.Fatalf("Failed to create s3 bucket: %v", err)
		}
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
