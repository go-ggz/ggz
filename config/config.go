package config

import (
	"time"
)

type database struct {
	Driver   string
	Username string
	Password string
	Name     string
	Host     string
	SSLMode  string
	Path     string
	TimeOut  int
}

type server struct {
	Host          string
	Addr          string
	Cert          string
	Key           string
	Root          string
	Storage       string
	Assets        string
	LetsEncrypt   bool
	StrictCurves  bool
	StrictCiphers bool
	Pprof         bool
	Token         string
	ShortenHost   string
	ShortenSize   int
	Cache         string
	Debug         bool `default:"true"`
}

type storage struct {
	Driver string
	Path   string
}

type admin struct {
	Users  []string
	Create bool
}

type cache struct {
	Driver string
	Expire int
	Prefix string
}

type session struct {
	Expire time.Duration
}

type qrcode struct {
	Enable bool
	Bucket string
}

type s3 struct {
	AccessID  string
	SecretKey string
	EndPoint  string
	SSL       bool
	Bucket    string
	Region    string
}

type auth0 struct {
	Key     string
	PemPath string
	Debug   bool
}

type prometheus struct {
	AuthToken string
}

type logs struct {
	Color  bool
	Debug  bool
	Pretty bool
}

// ContextKey for context package
type ContextKey string

func (c ContextKey) String() string {
	return "user context key " + string(c)
}

var (
	// Database represents the current database connection details.
	Database = &database{}

	// Server represents the informations about the server bindings.
	Server = &server{}

	// Admin represents the informations about the admin config.
	Admin = &admin{}

	// Session represents the informations about the session handling.
	Session = &session{}

	// Storage represents the informations about the storage bindings.
	Storage = &storage{}

	// QRCode represents the informations about the qrcode settings.
	QRCode = &qrcode{}

	// Minio represents the informations about the Minio server.
	Minio = &s3{}

	// Auth0 token information
	Auth0 = &auth0{}

	// ContextKeyUser for user
	ContextKeyUser = ContextKey("user")

	// Cache for redis, lur or memory cache
	Cache = &cache{}

	// Prometheus config
	Prometheus = &prometheus{}

	// Logs for zerolog
	Logs = &logs{}
)
