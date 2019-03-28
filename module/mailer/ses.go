package mailer

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/rs/zerolog/log"
)

const (
	// CharSet The character encoding for the email.
	CharSet = "UTF-8"
)

// SES for aws ses
type SES struct {
	sess    *session.Session
	source  *string
	to      []*string
	cc      []*string
	subject *string
	body    *string
}

func (c SES) From(name, address string) Mail {
	c.source = aws.String(fmt.Sprintf("%s <%s>", name, address))

	return c
}

func (c SES) To(address ...string) Mail {
	for _, v := range address {
		c.to = append(c.to, aws.String(v))
	}

	return c
}

func (c SES) Cc(address ...string) Mail {
	for _, v := range address {
		c.cc = append(c.cc, aws.String(v))
	}

	return c
}

func (c SES) Subject(subject string) Mail {
	c.subject = aws.String(subject)

	return c
}

func (c SES) Body(body string) Mail {
	c.body = aws.String(body)

	return c
}

// Send single email
func (c SES) Send() (interface{}, error) {
	// Create an SES session.
	svc := ses.New(c.sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: c.cc,
			ToAddresses: c.to,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    c.body,
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    c.body,
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    c.subject,
			},
		},
		Source: c.source,
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	resp, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Error().Err(aerr).Msg(ses.ErrCodeMessageRejected)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Error().Err(aerr).Msg(ses.ErrCodeMailFromDomainNotVerifiedException)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Error().Err(aerr).Msg(ses.ErrCodeConfigurationSetDoesNotExistException)
			default:
				log.Error().Err(aerr).Msg("AWS SES Error")
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error().Err(aerr).Msg("Unknown Error")
		}

		return nil, err
	}

	return resp, nil
}

// SESEngine initial ses
func SESEngine() (*SES, error) {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		return nil, err
	}

	return &SES{
		sess: sess,
	}, nil
}
