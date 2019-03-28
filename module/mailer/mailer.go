package mailer

import (
	"github.com/go-ggz/ggz/config"

	"github.com/rs/zerolog/log"
)

// Mail for smtp or ses interface
type Mail interface {
	From(string, string) Mail
	Send(config.Meta) (interface{}, error)
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
}

// Client for mail interface
var Client Mail

// NewEngine return storage interface
func NewEngine(c Config) (mail Mail, err error) {
	switch config.MailService.Driver {
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
