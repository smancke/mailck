package main

import (
	"github.com/smancke/mailck"
	"github.com/tarent/lib-compose/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const applicationName = "mailckd"

func main() {
	config := ReadConfig()
	if err := logging.Set(config.LogLevel, config.TextLogging); err != nil {
		exit(nil, err)
		return // return here for unittesing
	}

	logShutdownEvent()

	logging.LifecycleStart(applicationName, config)

	checkFunc := func(checkEmail string) (result mailck.CheckResult, textMessage string, err error) {
		return mailck.Check(config.FromEmail, checkEmail)
	}
	handlerChain := logging.NewLogMiddleware(NewValidationHandler(checkFunc))

	exit(nil, http.ListenAndServe(config.HostPort(), handlerChain))
}

func logShutdownEvent() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		exit(<-c, nil)
	}()
}

func exit(signal os.Signal, err error) {
	logging.LifecycleStop(applicationName, signal, err)
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	osExit(exitCode)
}

var osExit = func(exitCode int) {
	os.Exit(exitCode)
}
