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
	strVal := os.Getenv(name)
	if strVal == "" {
		log.WithFields(logrus.Fields{
			"name":    name,
			"default": def,
		}).Info("No ENV value, using default")
		return def
	}
	intVal, err := strconv.Atoi(strVal)
	if err != nil {
		log.WithFields(logrus.Fields{
			"name":    name,
			"value":   strVal,
			"default": def,
		}).Error("Invalid value, using default")
		return def
	}
	log.WithFields(logrus.Fields{
		"name":    name,
		"default": def,
		"value":   intVal,
	}).Info("Using ENV value")
	return int64(intVal)
}

// String provides a wrapper around os.Getenv with an optional fallback value.
func String(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	name := parts[0]
	def := ""
	if len(parts) > 1 {
		def = parts[1]
	}

	val := os.Getenv(name)
	if val == "" {
		log.WithFields(logrus.Fields{
			"name":    name,
			"default": def,
		}).Info("No ENV value, using default")
		return def
	}
	log.WithFields(logrus.Fields{
		"name":    name,
		"default": def,
		"value":   val,
	}).Info("Using ENV value")
	return val
}
