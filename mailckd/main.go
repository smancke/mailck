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

var exit = func(signal os.Signal, err error) {
	logging.LifecycleStop(applicationName, signal, err)
	if err == nil {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
