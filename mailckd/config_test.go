package main

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_ReadConfigDefaults(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"mailckd"}
	defer func() { os.Args = originalArgs }()

	d := DefaultConfig()
	assert.Equal(t, &d, ReadConfig())
}

func TestConfig_ReadConfig(t *testing.T) {
	input := []string{
		"--host=host",
		"--port=port",
		"--log-level=loglevel",
		"--text-logging=true",
		"--from-email=foo@example.com",
	}

	expected := &Config{
		Host:        "host",
		Port:        "port",
		LogLevel:    "loglevel",
		TextLogging: true,
		FromEmail:   "foo@example.com",
	}

	cfg, err := readConfig(flag.NewFlagSet("", flag.ContinueOnError), input)
	assert.NoError(t, err)
	assert.Equal(t, expected, cfg)
}

func TestConfig_ReadConfigFromEnv(t *testing.T) {
	assert.NoError(t, os.Setenv("MAILCKD_HOST", "host"))
	defer os.Unsetenv("MAILCKD_HOST")
	assert.NoError(t, os.Setenv("MAILCKD_PORT", "port"))
	defer os.Unsetenv("MAILCKD_PORT")
	assert.NoError(t, os.Setenv("MAILCKD_LOG_LEVEL", "loglevel"))
	defer os.Unsetenv("MAILCKD_LOG_LEVEL")
	assert.NoError(t, os.Setenv("MAILCKD_TEXT_LOGGING", "true"))
	defer os.Unsetenv("MAILCKD_TEXT_LOGGING")
	assert.NoError(t, os.Setenv("MAILCKD_FROM_EMAIL", "foo@example.com"))
	defer os.Unsetenv("MAILCKD_FROM_EMAIL")

	expected := &Config{
		Host:        "host",
		Port:        "port",
		LogLevel:    "loglevel",
		TextLogging: true,
		FromEmail:   "foo@example.com",
	}

	cfg, err := readConfig(flag.NewFlagSet("", flag.ContinueOnError), []string{})
	assert.NoError(t, err)
	assert.Equal(t, expected, cfg)
}
