package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/router"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/sync/errgroup"
	"gopkg.in/urfave/cli.v2"
)

var (
	defaultAddr        = ":8080"
	defaultShortenAddr = ":8081"
)

// Server provides the sub-command to start the API server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the gzz service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "db-driver",
				Value:       "mysql",
				Usage:       "Database driver selection",
				EnvVars:     []string{"GGZ_DB_DRIVER"},
				Destination: &config.Database.Driver,
			},
			&cli.StringFlag{
				Name:        "db-name",
				Value:       "test",
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
				Name:        "host",
				Value:       "http://localhost:8080",
				Usage:       "External access to server",
				EnvVars:     []string{"GGZ_SERVER_HOST"},
				Destination: &config.Server.Host,
			},
			&cli.StringFlag{
				Name:        "addr",
				Value:       defaultAddr,
				Usage:       "Address to bind the server",
				EnvVars:     []string{"GGZ_SERVER_ADDR"},
				Destination: &config.Server.Addr,
			},
			&cli.StringFlag{
				Name:        "shorten-addr",
				Value:       defaultShortenAddr,
				Usage:       "Address to bind the shorten server",
				EnvVars:     []string{"GGZ_SHORTEN_SERVER_ADDR"},
				Destination: &config.Server.ShortenAddr,
			},
			&cli.StringFlag{
				Name:        "root",
				Value:       "/",
				Usage:       "Root folder of the app",
				EnvVars:     []string{"GGZ_SERVER_ROOT"},
				Destination: &config.Server.Root,
			},
			&cli.StringFlag{
				Name:        "storage",
				Value:       "storage/",
				Usage:       "Folder for storing uploads",
				EnvVars:     []string{"GGZ_SERVER_STORAGE"},
				Destination: &config.Server.Storage,
			},
			&cli.BoolFlag{
				Name:        "pprof",
				Value:       false,
				Usage:       "Enable pprof debugging server",
				EnvVars:     []string{"GGZ_SERVER_PPROF"},
				Destination: &config.Server.Pprof,
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
			&cli.StringFlag{
				Name:        "path",
				Value:       "data/ggz.db",
				Usage:       "sqlite path",
				EnvVars:     []string{"GGZ_SQLITE_PATH"},
				Destination: &config.Database.Path,
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
				Value:       "",
				Usage:       "shorten-host",
				EnvVars:     []string{"GGZ_SERVER_SHORTEN_HOST"},
				Destination: &config.Server.ShortenHost,
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
					if config.Server.Addr != defaultAddr {
						logrus.Infof("With Let's Encrypt bind port have been overwritten!")
					}

					parsed, err := url.Parse(config.Server.Host)

					if err != nil {
						logrus.Fatal("Failed to parse host name. %s", err)
					}

					certManager := autocert.Manager{
						Prompt:     autocert.AcceptTOS,
						HostPolicy: autocert.HostWhitelist(parsed.Host),
						Cache:      autocert.DirCache(path.Join(config.Server.Storage, "certs")),
					}

					cfg.GetCertificate = certManager.GetCertificate

					var (
						g errgroup.Group
					)

					splitAddr := strings.SplitN(config.Server.Addr, ":", 2)
					logrus.Infof("Starting on %s:80 and %s:443", splitAddr[0], splitAddr[0])

					// load database
					router.GlobalInit()

					g.Go(func() error {
						return http.ListenAndServe(
							fmt.Sprintf("%s:80", splitAddr[0]),
							http.HandlerFunc(redirect),
						)
					})

					g.Go(func() error {
						return startServer(&http.Server{
							Addr:         fmt.Sprintf("%s:443", splitAddr[0]),
							Handler:      router.Load(),
							ReadTimeout:  5 * time.Second,
							WriteTimeout: 10 * time.Second,
							TLSConfig:    cfg,
						})
					})

					if err := g.Wait(); err != nil {
						logrus.Fatal(err)
					}
				} else {
					cert, err := tls.LoadX509KeyPair(
						config.Server.Cert,
						config.Server.Key,
					)

					if err != nil {
						logrus.Fatal("Failed to load SSL certificates. %s", err)
					}

					cfg.Certificates = []tls.Certificate{
						cert,
					}

					// load database
					router.GlobalInit()

					server := &http.Server{
						Addr:         config.Server.Addr,
						Handler:      router.Load(),
						ReadTimeout:  5 * time.Second,
						WriteTimeout: 10 * time.Second,
						TLSConfig:    cfg,
					}

					if err := startServer(server); err != nil {
						logrus.Fatal(err)
					}
				}
			} else {
				var (
					g errgroup.Group
				)

				// load global script
				logrus.Info("Initial project script.")
				router.GlobalInit()

				server01 := &http.Server{
					Addr:         config.Server.Addr,
					Handler:      router.Load(),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				server02 := &http.Server{
					Addr:         config.Server.ShortenAddr,
					Handler:      router.LoadRedirct(),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
				}

				g.Go(func() error {
					logrus.Infof("Starting main server on %s", config.Server.Addr)
					return startServer(server01)
				})

				g.Go(func() error {
					logrus.Infof("Starting shorten server on %s", config.Server.ShortenAddr)
					return startServer(server02)
				})

				if err := g.Wait(); err != nil {
					logrus.Fatal(err)
				}
			}

			return nil
		},
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	logrus.Debugf("Redirecting to %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}
