package forwarder

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"

	"github.com/gabihodoroaga/smtpd-email-forward/config"
	"github.com/gabihodoroaga/smtpd-email-forward/logger"
	"github.com/mailgun/mailgun-go/v4"
)

type mailgunForwarder struct {
	Domain string
	APIKey string
	URL string
}

func initMailGun() {
	if config.Config.MailgunDomain != "" && config.Config.MailgunAPIKey != "" {
		logger.Log.Info("Init Mailgun forwarder")
		forwarders = append(forwarders, &mailgunForwarder{
			Domain: config.Config.MailgunDomain,
			APIKey: config.Config.MailgunAPIKey,
			URL: config.Config.MailgunURL,
		})
	}
}

func (f mailgunForwarder) ForwardEmail(data []byte, recipient string) error {
	// create the client
	mg := mailgun.NewMailgun(f.Domain, f.APIKey)
	mg.SetAPIBase(f.URL)

	// create the message
	message := mg.NewMIMEMessage(ioutil.NopCloser(bytes.NewReader(data)), recipient)
	// create a context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		return err
	}

	logger.Log.Infof("Email sent successfully to mailgun: ID: %s Resp: %s", id, resp)

	return nil
}