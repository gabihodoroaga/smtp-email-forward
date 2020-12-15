package main

import (
	"github.com/gabihodoroaga/smtpd-email-forward/config"
	"github.com/gabihodoroaga/smtpd-email-forward/logger"
	"github.com/gabihodoroaga/smtpd-email-forward/smtpd"
	"github.com/gabihodoroaga/smtpd-email-forward/forwarder"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("smtpd server is starting...")
	config.InitConfig()
	forwarder.InitForwarders()
	smtpd.StartServer()
}
