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
	ShortenAddr   string
	ShortenHost   string
}

type storage struct {
	Driver string
	Path   string
}

type admin struct {
	Users  []string
	Create bool
}

type session struct {
	Expire time.Duration
}

type qrcode struct {
	Enable bool
	Bucket string
}

var (
	// Debug represents the flag to enable or disable debug logging.
	Debug bool

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
)
