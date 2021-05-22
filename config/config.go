package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/gabihodoroaga/smtpd-email-forward/logger"
	"github.com/ilyakaznacheev/cleanenv"
)

type args struct {
	ConfigPath string
}

// Configuration struct ...
type Configuration struct {
	AppName        string   `yaml:"appname" env:"APP_NAME"`
	Hostname       string   `yaml:"hostname" env:"HOST_NAME"`
	Domains        []string `yaml:"domains" env:"DOMAINS" env-separator:","`
	StartTLSCert   string   `yaml:"starttlscert" env:"START_TLS_CERT" env-default:"certs/server.crt"`
	StartTLSKey    string   `yaml:"starttlskey" env:"START_TLS_KEY" env-default:"certs/server.key"`
	GCPBucket      string   `yaml:"gcpbucket" env:"GCP_BUCKET"`
	GCPCredentials string   `yaml:"gcpcredentials" env:"GCP_CREDENTIALS" env-default:"certs/gcpServiceAccount.json"`
	ForwardTo      string   `yaml:"forwardto" env:"FORWARD_TO"`
	MailgunDomain  string   `yaml:"mailgundomain" env:"MAILGUN_DOMAIN"`
	MailgunAPIKey  string   `yaml:"mailgunapikey" env:"MAILGUN_API_KEY"`
	MailgunURL     string   `yaml:"mailgunurl" env:"MAILGUN_URL" env-default:"https://api.mailgun.net/v3"`
}

// Config is the application configuration
var Config Configuration

// InitConfig - Initialize the the application configuration
func InitConfig() {
	args := processArgs(&Config)
	// read configuration from the file and environment variables
	if err := cleanenv.ReadConfig(args.ConfigPath, &Config); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// TODO: print configuration values
	logger.Log.Info("Configuration loaded successfully")
}

// ProcessArgs processes and handles CLI arguments
func processArgs(cfg interface{}) args {
	var a args

	f := flag.NewFlagSet("Example server", 1)
	f.StringVar(&a.ConfigPath, "c", "config.yml", "Path to configuration file")

	fu := f.Usage
	f.Usage = func() {
		fu()
		envHelp, _ := cleanenv.GetDescription(cfg, nil)
		fmt.Fprintln(f.Output())
		fmt.Fprintln(f.Output(), envHelp)
	}

	f.Parse(os.Args[1:])
	return a
}
