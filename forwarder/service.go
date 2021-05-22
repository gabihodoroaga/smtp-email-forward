package forwarder

import (
	"errors"
	"strings"

	"github.com/gabihodoroaga/smtpd-email-forward/config"
)

var forwarders []mailForwarder

type mailForwarder interface {
	ForwardEmail(data []byte, recipient string) error
}

// InitForwarders initialize all the email forwarders
func InitForwarders() {
	initMailGun()
}

// ForwardEmail forwards the email to all the configured forwarders
func ForwardEmail(data []byte) error {
	var errstrings []string
	for _, f := range forwarders {
		if err := f.ForwardEmail(data, config.Config.ForwardTo); err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}
	if len(errstrings) > 0 {
		return errors.New(strings.Join(errstrings, "\n"))
	}
	return nil
}
