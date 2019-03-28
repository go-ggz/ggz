package main

import (
	"fmt"
	"sync"

	"github.com/go-ggz/ggz/config"
	"github.com/go-ggz/ggz/module/mailer"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/rs/zerolog/log"
	"gopkg.in/urfave/cli.v2"
)

// Mail provides the sub-command to send email.
func Mail() *cli.Command {
	return &cli.Command{
		Name:  "email",
		Usage: "send test email",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "aws-access-id",
				Usage:       "aws access key",
				EnvVars:     []string{"AWS_ACCESS_KEY_ID"},
				Destination: &config.AWS.AccessID,
			},
			&cli.StringFlag{
				Name:        "aws-secret-key",
				Usage:       "aws secret key",
				EnvVars:     []string{"AWS_SECRET_ACCESS_KEY"},
				Destination: &config.AWS.SecretKey,
			},
			&cli.StringFlag{
				Name:        "mail-driver",
				Usage:       "mail driver",
				Value:       "ses",
				EnvVars:     []string{"GGZ_MAIL_DRIVER"},
				Destination: &config.MailService.Driver,
			},
			&cli.StringFlag{
				Name:        "smtp-host",
				Usage:       "smtp host",
				EnvVars:     []string{"GGZ_SMTP_HOST"},
				Destination: &config.SMTP.Host,
			},
			&cli.StringFlag{
				Name:        "smtp-port",
				Usage:       "smtp port",
				EnvVars:     []string{"GGZ_SMTP_PORT"},
				Destination: &config.SMTP.Port,
			},
			&cli.StringFlag{
				Name:        "smtp-username",
				Usage:       "smtp username",
				EnvVars:     []string{"GGZ_SMTP_USERNAME"},
				Destination: &config.SMTP.Username,
			},
			&cli.StringFlag{
				Name:        "smtp-password",
				Usage:       "smtp password",
				EnvVars:     []string{"GGZ_SMTP_PASSWORD"},
				Destination: &config.SMTP.Password,
			},
			&cli.IntFlag{
				Name:  "email-count",
				Usage: "send email count",
				Value: 1,
			},
		},
		Action: func(c *cli.Context) error {
			// initial mailer service
			if _, err := mailer.NewEngine(mailer.Config{}); err != nil {
				log.Fatal().Err(err).Msgf("failed to initial mailer")
			}

			count := c.Int("email-count")
			wg := sync.WaitGroup{}

			for i := 0; i < count; i++ {
				wg.Add(1)
				go func(k int, wg *sync.WaitGroup) {
					data := config.Meta{
						Sender: config.To{
							Name:  "Bo-Yi Wu",
							Email: "appleboy.tw@gmail.com",
						},
						Subject:     fmt.Sprintf("[Fullinn] 測試 AWS 電子郵件系統 [%d]", (k + 1)),
						ToAddresses: []string{"appleboy.tw@gmail.com"},
						Body: "<h1>繁體中文 Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
							"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
							"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>",
					}

					resp, err := mailer.Client.Send(data)

					if err != nil {
						log.Fatal().Err(err).Msgf("failed to send email")
					}

					log.Info().Msgf("Send the Email completely [%d]", (k + 1))
					if v, ok := resp.(*ses.SendEmailOutput); ok {
						log.Info().Msgf("Message ID: %s", *v.MessageId)
					}
					wg.Done()
				}(i, &wg)
			}

			wg.Wait()

			return nil
		},
	}
}
