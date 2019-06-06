package mailer

import (
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
)

type from struct {
	Name    string
	Address string
}

// SMTP setting
type SMTP struct {
	host     string
	port     string
	username string
	password string
	from     from
	to       []string
	cc       []string
	subject  string
	body     string
}

// From for sender information
func (c SMTP) From(name, address string) Mail {
	c.from = from{
		Name:    name,
		Address: address,
	}

	return c
}

// To for mailto list
func (c SMTP) To(address ...string) Mail {
	c.to = address

	return c
}

// Cc for cc list
func (c SMTP) Cc(address ...string) Mail {
	c.cc = address

	return c
}

// Subject for email title
func (c SMTP) Subject(subject string) Mail {
	c.subject = subject

	return c
}

// Body for email body
func (c SMTP) Body(body string) Mail {
	c.body = body

	return c
}

// Send email
func (c SMTP) Send() (interface{}, error) {
	m := email.NewHTMLMessage(c.subject, c.body)
	m.From = mail.Address{
		Name:    c.from.Name,
		Address: c.from.Address,
	}
	m.To = c.to
	m.Cc = c.cc

	// send it
	auth := smtp.PlainAuth("", c.username, c.password, c.host)

	return nil, email.Send(c.host+":"+c.port, auth, m)
}

// SMTPEngine initial smtp object
func SMTPEngine(host, port, username, password string) (*SMTP, error) {
	return &SMTP{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}, nil
}
