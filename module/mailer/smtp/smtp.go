package smtp

import (
	"net/mail"
	"net/smtp"

	"github.com/go-ggz/ggz/config"

	"github.com/scorredoira/email"
)

// Client for smtp
type Client struct {
	host     string
	port     string
	username string
	password string
}

// Send single email
func (c *Client) Send(meta config.Meta) (interface{}, error) {
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
func NewEngine(host, port, username, password string) (*Client, error) {
	return &Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}, nil
}
