package main

import (
	"flag"
	"github.com/caarlos0/env"
	"os"
)

func DefaultConfig() Config {
	return Config{
		Host:      "localhost",
		Port:      "6788",
		LogLevel:  "info",
		FromEmail: "noreply@mailck.io",
	}
}

type Config struct {
	Host        string `env:"MAILCKD_HOST"`
	Port        string `env:"MAILCKD_PORT"`
	LogLevel    string `env:"MAILCKD_LOG_LEVEL"`
	TextLogging bool   `env:"MAILCKD_TEXT_LOGGING"`
	FromEmail   string `env:"MAILCKD_FROM_EMAIL"`
}

func (c Config) HostPort() string {
	return c.Host + ":" + c.Port
}

func ReadConfig() *Config {
	c, err := readConfig(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:])
	if err != nil {
		// should never happen, because of flag default policy ExitOnError
		panic(err)
	}
	return c
}

func readConfig(f *flag.FlagSet, args []string) (*Config, error) {
	config := DefaultConfig()

	// Environment variables
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}

	f.StringVar(&config.Host, "host", config.Host, "The host to listen on")
	f.StringVar(&config.Port, "port", config.Port, "The port to listen on")
	f.StringVar(&config.LogLevel, "log-level", config.LogLevel, "The log level")
	f.BoolVar(&config.TextLogging, "text-logging", config.TextLogging, "Log in text format instead of json")
	f.StringVar(&config.FromEmail, "from-email", config.FromEmail, "The from email when connecting to the mailserver")

	// Arguments variables
	err = f.Parse(args)
	if err != nil {
		return nil, err
	}

	return &config, err
}
