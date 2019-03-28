package ses

import (
	"fmt"

	"github.com/go-ggz/ggz/config"

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

// Client for ses
type Client struct {
	sess *session.Session
}

// Send single email
func (c *Client) Send(meta config.Meta) (interface{}, error) {
	toAddresses := []*string{}
	ccAddresses := []*string{}
	for _, v := range meta.ToAddresses {
		toAddresses = append(toAddresses, aws.String(v))
	}

	for _, v := range meta.CcAddresses {
		ccAddresses = append(ccAddresses, aws.String(v))
	}

	// Create an SES session.
	svc := ses.New(c.sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: ccAddresses,
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(meta.Body),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(meta.Body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(meta.Subject),
			},
		},
		Source: aws.String(fmt.Sprintf("%s <%s>", meta.Sender.Name, meta.Sender.Email)),
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

// NewEngine initial ses
func NewEngine() (*Client, error) {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		return nil, err
	}

	return &Client{
		sess: sess,
	}, nil
}
