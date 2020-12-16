package smtpd

import (
	"context"
	"net"
	"strings"

	"github.com/gabihodoroaga/smtpd-email-forward/config"
	"github.com/gabihodoroaga/smtpd-email-forward/logger"
	"github.com/gabihodoroaga/smtpd-email-forward/storage"
	"github.com/gabihodoroaga/smtpd-email-forward/forwarder"
	"github.com/mhale/smtpd"
)

// StartServer initialize and start the smtp server on
func StartServer() {
	addr := "0.0.0.0:2525"
	logger.Log.Infof("Start listening on: %s", addr)
	err := listenAndServeTLS(
		addr,
		config.Config.StartTLSCert,
		config.Config.StartTLSKey,
		mailHandler,
		rcptHandler,
		config.Config.AppName,
		config.Config.Hostname)
	if err != nil {
		panic(err)
	}
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	logger.Log.Infof("Received mail from %s for %s", from, strings.Join(to, ";"))
	// save the email to GCP bucket
	if config.Config.GCPBucket != "" {
		ctx := context.Background()
		for _, m := range to {
			filename, err := storage.UploadFileToBucket(ctx, m, data)
			if err != nil {
				logger.Log.Errorf("Error saving file %s to gcs: %s", filename, err)
			} else {
				logger.Log.Infof("Save mail from %s for %s to file %s", from, to, filename)
			}
		}
	}
	// forward the emails
	err := forwarder.ForwardEmail(data)

	if err != nil {
		logger.Log.Errorf("Errors forward email:%s", err)
	}
}

func rcptHandler(remoteAddr net.Addr, from string, to string) bool {
	domain := getDomain(to)
	if domain == "" {
		return false
	}
	for _, b := range config.Config.Domains {
		if b == domain {
			return true
		}
	}
	return false
}

func getDomain(email string) string {
	at := strings.LastIndex(email, "@")
	if at >= 0 {
		return email[at+1:]
	}
	return ""
}

func listenAndServeTLS(addr string, certFile string, keyFile string, handler smtpd.Handler, rcpt smtpd.HandlerRcpt, appname string, hostname string) error {
	srv := &smtpd.Server{
		Addr:        addr,
		Handler:     handler,
		HandlerRcpt: rcpt,
		Appname:     appname,
		Hostname:    hostname,
		TLSRequired: false}
	err := srv.ConfigureTLS(certFile, keyFile)
	if err != nil {
		return err
	}
	return srv.ListenAndServe()
}
