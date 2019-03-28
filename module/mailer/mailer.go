package mailer

import (
	"github.com/rs/zerolog/log"
)

// Mail for smtp or ses interface
type Mail interface {
	From(string, string) Mail
	To(...string) Mail
	Cc(...string) Mail
	Subject(string) Mail
	Body(string) Mail
	Send() (interface{}, error)
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Driver   string
}

// Client for mail interface
var Client Mail

// NewEngine return storage interface
func NewEngine(c Config) (mail Mail, err error) {
	switch c.Driver {
	case "smtp":
		Client, err = SMTPEngine(
			c.Host,
			c.Port,
			c.Username,
			c.Password,
		)
		if err != nil {
			return nil, err
		}
	case "ses":
		Client, err = SESEngine()
		if err != nil {
			return nil, err
		}
	default:
		log.Error().Msg("Unknown email driver")
	}

	return mail, nil
}
