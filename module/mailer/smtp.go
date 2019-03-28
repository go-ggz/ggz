package mailer

import (
	"net/mail"
	"net/smtp"

	"github.com/go-ggz/ggz/config"

	"github.com/scorredoira/email"
)

type from struct {
	Name    string
	Address string
}

// Client for smtp
type SMTP struct {
	host     string
	port     string
	username string
	password string
	from     from
}

func (c SMTP) From(name, address string) Mail {
	c.from = from{
		Name:    name,
		Address: address,
	}

	return c
}

// Send single email
func (c SMTP) Send(meta config.Meta) (interface{}, error) {
	m := email.NewHTMLMessage(meta.Subject, meta.Body)
	m.From = mail.Address{
		Name:    meta.Sender.Name,
		Address: meta.Sender.Email,
	}
	m.To = meta.ToAddresses
	m.Cc = meta.CcAddresses

	// send it
	auth := smtp.PlainAuth("", c.username, c.password, c.host)

	return nil, email.Send(c.host+":"+c.port, auth, m)
}

// NewEngine initial smtp
func SMTPEngine(host, port, username, password string) (*SMTP, error) {
	return &SMTP{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}, nil
}
