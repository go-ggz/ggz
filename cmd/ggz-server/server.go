package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/router"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
	"gopkg.in/urfave/cli.v2"
)

var (
	defaultHostAddr = ":8080"
)

// Server provides the sub-command to start the API server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the gzz service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "assets",
				Value:       "",
				Usage:       "Path to custom assets and templates",
				EnvVars:     []string{"GGZ_SERVER_ASSETS"},
				Destination: &config.Server.Assets,
			},
			&cli.StringFlag{
				Name:        "db-driver",
				Value:       "sqlite3",
				Usage:       "Database driver selection",
				EnvVars:     []string{"GGZ_DB_DRIVER"},
				Destination: &config.Database.Driver,
			},
			&cli.StringFlag{
				Name:        "db-name",
				Value:       "ggz",
				Usage:       "Name for database connection",
				EnvVars:     []string{"GGZ_DB_NAME"},
				Destination: &config.Database.Name,
			},
			&cli.StringFlag{
				Name:        "db-username",
				Value:       "root",
				Usage:       "Username for database connection",
				EnvVars:     []string{"GGZ_DB_USERNAME"},
				Destination: &config.Database.Username,
			},
			&cli.StringFlag{
				Name:        "db-password",
				Value:       "root",
				Usage:       "Password for database connection",
				EnvVars:     []string{"GGZ_DB_PASSWORD"},
				Destination: &config.Database.Password,
			},
			&cli.StringFlag{
				Name:        "db-host",
				Value:       "localhost:3306",
				Usage:       "Host for database connection",
				EnvVars:     []string{"GGZ_DB_HOST"},
				Destination: &config.Database.Host,
			},
			&cli.StringFlag{
				Name:        "path",
				Value:       "data/db/ggz.db",
				Usage:       "sqlite path",
				EnvVars:     []string{"GGZ_SQLITE_PATH"},
				Destination: &config.Database.Path,
			},
			&cli.StringFlag{
				Name:        "host",
				Value:       "http://localhost:8080",
				Usage:       "External access to server",
				EnvVars:     []string{"GGZ_SERVER_HOST"},
				Destination: &config.Server.Host,
			},
			&cli.StringFlag{
				Name:        "addr",
				Value:       defaultHostAddr,
				Usage:       "Address to bind the server",
				EnvVars:     []string{"GGZ_SERVER_ADDR"},
				Destination: &config.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "root",
				Value:       "/",
				Usage:       "Root folder of the app",
				EnvVars:     []string{"GGZ_SERVER_ROOT"},
				Destination: &config.Server.Root,
			},
			&cli.BoolFlag{
				Name:        "pprof",
				Value:       false,
				Usage:       "Enable pprof debugging server",
				EnvVars:     []string{"GGZ_SERVER_PPROF"},
				Destination: &config.Server.Pprof,
			},
			&cli.BoolFlag{
				Name:        "graphiql",
				Value:       false,
				Usage:       "Enable graphiql interface",
				EnvVars:     []string{"GOBENTO_SERVER_GRAPHIQL"},
				Destination: &config.Server.GraphiQL,
			},
			&cli.StringFlag{
				Name:        "cert",
				Value:       "",
				Usage:       "Path to SSL cert",
				EnvVars:     []string{"GGZ_SERVER_CERT"},
				Destination: &config.Server.Cert,
			},
			&cli.StringFlag{
				Name:        "key",
				Value:       "",
				Usage:       "Path to SSL key",
				EnvVars:     []string{"GGZ_SERVER_KEY"},
				Destination: &config.Server.Key,
			},
			&cli.BoolFlag{
				Name:        "letsencrypt",
				Value:       false,
				Usage:       "Enable Let's Encrypt SSL",
				EnvVars:     []string{"GGZ_SERVER_LETSENCRYPT"},
				Destination: &config.Server.LetsEncrypt,
			},
			&cli.BoolFlag{
				Name:        "strict-curves",
				Value:       false,
				Usage:       "Use strict SSL curves",
				EnvVars:     []string{"GGZ_STRICT_CURVES"},
				Destination: &config.Server.StrictCurves,
			},
			&cli.BoolFlag{
				Name:        "strict-ciphers",
				Value:       false,
				Usage:       "Use strict SSL ciphers",
				EnvVars:     []string{"GGZ_STRICT_CIPHERS"},
				Destination: &config.Server.StrictCiphers,
			},
			&cli.DurationFlag{
				Name:        "expire",
				Value:       time.Hour * 24,
				Usage:       "Session expire duration",
				EnvVars:     []string{"GGZ_SESSION_EXPIRE"},
				Destination: &config.Session.Expire,
			},
			&cli.StringSliceFlag{
				Name:    "admin-user",
				Value:   &cli.StringSlice{},
				Usage:   "Enforce user as an admin",
				EnvVars: []string{"GGZ_ADMIN_USERS"},
			},
			&cli.BoolFlag{
				Name:        "admin-create",
				Value:       true,
				Usage:       "Create an initial admin user",
				EnvVars:     []string{"GGZ_ADMIN_CREATE"},
				Destination: &config.Admin.Create,
			},
			&cli.StringFlag{
				Name:        "token",
				Value:       "",
				Usage:       "Header token",
				EnvVars:     []string{"GGZ_TOKEN"},
				Destination: &config.Server.Token,
			},
			&cli.IntFlag{
				Name:        "timeout",
				Value:       500,
				Usage:       "sqlite database timeout",
				EnvVars:     []string{"GGZ_SQLITE_TIMEOUT"},
				Destination: &config.Database.TimeOut,
			},
			&cli.StringFlag{
				Name:        "mode",
				Value:       "",
				Usage:       "databas ssl mode",
				EnvVars:     []string{"GGZ_SSL_MODE"},
				Destination: &config.Database.SSLMode,
			},
			&cli.StringFlag{
				Name:        "shorten-host",
				Value:       "http://localhost:8081",
				Usage:       "shorten-host",
				EnvVars:     []string{"GGZ_SERVER_SHORTEN_HOST"},
				Destination: &config.Server.ShortenHost,
			},
			&cli.IntFlag{
				Name:        "shorten-size",
				Value:       5,
				Usage:       "shorten-size",
				EnvVars:     []string{"GGZ_SERVER_SHORTEN_SIZE"},
				Destination: &config.Server.ShortenSize,
			},
			&cli.StringFlag{
				Name:        "storage-driver",
				Value:       "disk",
				Usage:       "Storage driver selection",
				EnvVars:     []string{"GGZ_STORAGE_DRIVER"},
				Destination: &config.Storage.Driver,
			},
			&cli.StringFlag{
				Name:        "storage-path",
				Value:       "storage/",
				Usage:       "Folder for storing uploads",
				EnvVars:     []string{"GGZ_STORAGE_PATH"},
				Destination: &config.Storage.Path,
			},
			&cli.BoolFlag{
				Name:        "qrcode-enable",
				Usage:       "qrcode module enable",
				EnvVars:     []string{"GGZ_QRCODE_ENABLE"},
				Destination: &config.QRCode.Enable,
			},
			&cli.StringFlag{
				Name:        "qrcode-bucket",
				Value:       "qrcode",
				Usage:       "qrcode bucket name",
				EnvVars:     []string{"GGZ_QRCODE_BUCKET"},
				Destination: &config.QRCode.Bucket,
			},
			&cli.StringFlag{
				Name:        "minio-access-id",
				Value:       "",
				Usage:       "minio-access-id",
				EnvVars:     []string{"GGZ_MINIO_ACCESS_ID"},
				Destination: &config.Minio.AccessID,
			},
			&cli.StringFlag{
				Name:        "minio-secret-key",
				Value:       "",
				Usage:       "minio-secret-key",
				EnvVars:     []string{"GGZ_MINIO_SECRET_KEY"},
				Destination: &config.Minio.SecretKey,
			},
			&cli.StringFlag{
				Name:        "minio-endpoint",
				Value:       "",
				Usage:       "minio-endpoint",
				EnvVars:     []string{"GGZ_MINIO_ENDPOINT"},
				Destination: &config.Minio.EndPoint,
			},
			&cli.BoolFlag{
				Name:        "minio-ssl",
				Usage:       "minio-ssl",
				EnvVars:     []string{"GGZ_MINIO_SSL"},
				Destination: &config.Minio.SSL,
			},
			&cli.StringFlag{
				Name:        "minio-bucket",
				Value:       "qrcode",
				Usage:       "minio-bucket",
				EnvVars:     []string{"GGZ_MINIO_BUCKET"},
				Destination: &config.Minio.Bucket,
			},
			&cli.StringFlag{
				Name:        "minio-region",
				Value:       "us-east-1",
				Usage:       "minio-region",
				EnvVars:     []string{"GGZ_MINIO_REGION"},
				Destination: &config.Minio.Region,
			},
			&cli.StringFlag{
				Name:        "auth0-pem-path",
				Usage:       "Auth0 Pem file path",
				EnvVars:     []string{"GGZ_AUTH0_PEM_PATH"},
				Destination: &config.Auth0.PemPath,
			},
			&cli.BoolFlag{
				Name:        "auth0-debug",
				Usage:       "Auth0 debug",
				EnvVars:     []string{"GGZ_AUTH0_DEBUG"},
				Destination: &config.Auth0.Debug,
			},
			&cli.StringFlag{
				Name:        "auth0-key-name",
				Usage:       "Auth0 key content",
				EnvVars:     []string{"GGZ_AUTH0_Key"},
				Destination: &config.Auth0.Key,
			},
			&cli.StringFlag{
				Name:        "cache-driver",
				Value:       "default",
				Usage:       "Cache driver selection",
				EnvVars:     []string{"GGZ_CACHE_DRIVER"},
				Destination: &config.Cache.Driver,
			},
			&cli.IntFlag{
				Name:        "cache-expire-time",
				Value:       15,
				Usage:       "cache expire time (minutes)",
				EnvVars:     []string{"GGZ_CACHE_EXPIRE"},
				Destination: &config.Cache.Expire,
			},
			&cli.StringFlag{
				Name:        "cache-prefix-name",
				Value:       "ggz",
				Usage:       "prefix name of key",
				EnvVars:     []string{"GGZ_CACHE_PREFIX_NAME"},
				Destination: &config.Cache.Prefix,
			},
			&cli.StringFlag{
				Name:        "prometheus-auth-token",
				EnvVars:     []string{"GGZ_PROMETHEUS_AUTH_TOKEN"},
				Usage:       "token to secure prometheus metrics endpoint",
				Destination: &config.Prometheus.AuthToken,
			},
			&cli.StringFlag{
				Name:        "auth-driver",
				EnvVars:     []string{"GGZ_AUTH_DRIVER"},
				Usage:       "auth driver",
				Value:       "auth0",
				Destination: &config.Auth.Driver,
			},
		},
		Before: func(c *cli.Context) error {
			if len(c.StringSlice("admin-user")) > 0 {
				// StringSliceFlag doesn't support Destination
				config.Admin.Users = c.StringSlice("admin-user")
			}

			return nil
		},
		Action: func(c *cli.Context) error {
			idleConnsClosed := make(chan struct{})

			// load global script
			log.Info().Msg("Initial module engine.")
			router.GlobalInit()

			server := &http.Server{
				Addr:         config.Server.Addr,
				Handler:      router.Load(),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			go func(srv *http.Server) {
				sigint := make(chan os.Signal, 1)

				// interrupt signal sent from terminal
				signal.Notify(sigint, os.Interrupt)
				// sigterm signal sent from kubernetes
				signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
				defer signal.Stop(sigint)

				<-sigint

				log.Info().Msg("received an interrupt signal, shut down the server.")
				// We received an interrupt signal, shut down.
				if err := srv.Shutdown(context.Background()); err != nil {
					// Error from closing listeners, or context timeout:
					log.Error().Err(err).Msg("HTTP server Shutdown")
				}
				close(idleConnsClosed)
			}(server)

			if config.Server.LetsEncrypt || (config.Server.Cert != "" && config.Server.Key != "") {
				cfg := &tls.Config{
					PreferServerCipherSuites: true,
					MinVersion:               tls.VersionTLS12,
				}

				if config.Server.StrictCurves {
					cfg.CurvePreferences = []tls.CurveID{
						tls.CurveP521,
						tls.CurveP384,
						tls.CurveP256,
					}
				}

				if config.Server.StrictCiphers {
					cfg.CipherSuites = []uint16{
						tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
						tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
						tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					}
				}

				if config.Server.LetsEncrypt {
					if config.Server.Addr != defaultHostAddr {
						log.Fatal().Msg("With Let's Encrypt bind port have been overwritten!")
					}

					parsed, err := url.Parse(config.Server.Host)

					if err != nil {
						log.Fatal().Err(err).Msg("Failed to parse host name.")
					}

					certManager := &autocert.Manager{
						Prompt:     autocert.AcceptTOS,
						HostPolicy: autocert.HostWhitelist(parsed.Host),
						Cache:      autocert.DirCache(path.Join(config.Server.Storage, "certs")),
					}

					cfg.GetCertificate = certManager.GetCertificate

					var (
						g errgroup.Group
					)

					splitAddr := strings.SplitN(config.Server.Addr, ":", 2)
					log.Info().Msgf("Starting on %s:80 and %s:443", splitAddr[0], splitAddr[0])

					g.Go(func() error {
						return http.ListenAndServe(
							fmt.Sprintf("%s:80", splitAddr[0]),
							certManager.HTTPHandler(http.HandlerFunc(redirect)),
						)
					})

					g.Go(func() error {
						server.Addr = fmt.Sprintf("%s:443", splitAddr[0])
						server.TLSConfig = cfg
						return startServer(server)
					})

					if err := g.Wait(); err != nil {
						log.Fatal().Err(err)
					}
				} else {
					cert, err := tls.LoadX509KeyPair(
						config.Server.Cert,
						config.Server.Key,
					)

					if err != nil {
						log.Fatal().Err(err).Msg("Failed to load SSL certificates.")
					}

					cfg.Certificates = []tls.Certificate{
						cert,
					}

					// Add TLS config
					server.TLSConfig = cfg

					if err := startServer(server); err != nil {
						log.Fatal().Err(err)
					}
				}
			} else {
				var (
					g errgroup.Group
				)

				g.Go(func() error {
					log.Info().Msgf("Starting shorten server on %s", config.Server.Addr)
					return startServer(server)
				})

				if err := g.Wait(); err != nil {
					log.Fatal().Err(err)
				}
			}

			<-idleConnsClosed

			return nil
		},
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	log.Printf("Redirecting to %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

func startServer(s *http.Server) error {
	if s.TLSConfig == nil {
		return s.ListenAndServe()
	}
	return s.ListenAndServeTLS("", "")
}
