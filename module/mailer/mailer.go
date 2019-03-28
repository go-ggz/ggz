package mailer

import (
	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/module/mailer/ses"
	"github.com/go-ggz/ggz/module/mailer/smtp"

	"github.com/rs/zerolog/log"
)

// Mail for smtp or ses interface
type Mail interface {
	Send(config.Meta) (interface{}, error)
}

// Client for mail interface
var Client Mail

// NewEngine return storage interface
func NewEngine() (main Mail, err error) {
	switch config.MailService.Driver {
	case "smtp":
		Client, err = smtp.NewEngine()
		if err != nil {
			return nil, err
		}
	case "ses":
		Client, err = ses.NewEngine()
		if err != nil {
			return nil, err
		}
	default:
		log.Error().Msg("Unknown email driver")
	}

	return mail, nil
}
