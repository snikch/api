package log

import (
	"os"

	"github.com/sebest/logrusly"
	"github.com/sirupsen/logrus"
	"github.com/snikch/api/lifecycle"
)

// Logger is a globally shared logrus instance.
var Logger = logrus.New()

// Convenience access to the common logging functions.
var (
	WithError  = Logger.WithError
	WithFields = Logger.WithFields
	WithField  = Logger.WithField
	Debug      = Logger.Debug
	Info       = Logger.Info
	Warn       = Logger.Warn
	Error      = Logger.Error
	Fatal      = Logger.Fatal
)

func init() {

	switch os.Getenv("LOG_FORMAT") {
	// Log in json if LOG_FORMAT is set to json
	case "json":
		Logger.Formatter = &logrus.JSONFormatter{}
	}

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		Logger.Level = logrus.DebugLevel
	case "warning":
		Logger.Level = logrus.WarnLevel
	default:
		Logger.Level = logrus.InfoLevel
	}
}

// InitLoggly will create a loggly hook and register it. It will also register
// a lifecycle shutdown callback to flush any remaining logs, if you're using
// the lifecycle package.
func InitLoggly(key, applicationID string, tags ...string) {
	// No-op if no key is set.
	if key == "" {
		return
	}

	hook := logrusly.NewLogglyHook(
		key,
		applicationID,
		logrus.InfoLevel,
		tags...,
	)
	Logger.Hooks.Add(hook)
	lifecycle.RegisterShutdownCallback("loggly", func() error {
		hook.Flush()
		return nil
	})
}
