package config

import (
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/snikch/api/log"
)

// init is responsible for initializing the entire configuration environment.
// In development, this means loading in and environment variables from both
// .env.default and .env. The .env.default file is commited into source control
// to provide sane defaults without sharing any secrets, whereas .env is ignored
// and should be used to store local development secrets, such as API keys for
// any services required to run.
func init() {
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".env.default")
	baseDir := os.Getenv("BASE_DIR")
	if baseDir != "" {
		_ = godotenv.Load(baseDir + ".env")
		_ = godotenv.Load(baseDir + ".env.default")
	}

}

// Int provides a convenience wrapper for retrieving an integer from an env var.
// If the env var is not set, the supplied default will be returned.
func Int(name string, def int64) int64 {
	return int(name, def, false)
}

// PrivateInt provides a convenience wrapper for retrieving an integer from an env var.
// If the env var is not set, the supplied default will be returned. No values
// are logged from this function.
func PrivateInt(name string, def int64) int64 {
	return int(name, def, true)
}

func int(name string, def int64, private bool) int64 {
	strVal := os.Getenv(name)
	l := log.WithField("name", name)

	// No value? Use the default value provided.
	if strVal == "" {
		if !private {
			l = l.WithField("default", def)
		}
		l.Info("No ENV value, using default")
		return def
	}

	// Attempt conversion from string to an int.
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		if !private {
			l = log.WithFields(logrus.Fields{
				"value":   strVal,
				"default": def,
			})
		}
		l.Error("Invalid value, using default")
		return def
	}

	// If we're here, it's all good. Use the new value.
	if !private {
		l = log.WithFields(logrus.Fields{
			"default": def,
			"value":   intVal,
		})
	}
	l.Info("Using ENV value")
	return int64(intVal)
}

// String provides a wrapper around os.Getenv with an optional fallback value.
func String(parts ...string) string {
	return str(false, parts...)
}

// PrivateString provides a wrapper around os.Getenv with an optional fallback
// value. No values are logged from this function.
func PrivateString(parts ...string) string {
	return str(true, parts...)
}

func str(private bool, parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	name := parts[0]
	def := ""
	if len(parts) > 1 {
		def = parts[1]
	}

	l := log.WithField("name", name)

	val := os.Getenv(name)
	if val == "" {
		if !private {
			l = l.WithField("default", def)
		}
		l.Info("No ENV value, using default")
		return def
	}

	if !private {
		l = l.WithFields(logrus.Fields{
			"default": def,
			"value":   val,
		})
	}
	l.Info("Using ENV value")
	return val
}
