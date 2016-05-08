package dbtest

import (
	"database/sql"
	"os"
	"strings"

	"github.com/snikch/api/config"
	"github.com/snikch/api/fail"
	"github.com/snikch/api/log"
	"github.com/snikch/goose/lib/goose"
)

// SetUp creates the database.
func SetUp(seed string) {
	conf := newConf()

	// Retrieve the original DSN value
	originalDSN := conf.Driver.OpenStr

	// Remove the database name
	conf.Driver.OpenStr = strings.Replace(conf.Driver.OpenStr, "DBNAME", "", -1)
	db, err := sql.Open(conf.Driver.Name, conf.Driver.OpenStr)
	if err != nil {
		log.WithError(fail.Trace(err)).Fatal("Could not open db connection")
	}

	// Create a database for this test.
	dbName := "test_" + seed
	_, err = db.Query("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.WithError(fail.Trace(err)).WithField("name", dbName).Fatal("Could not create database")
	}

	// Replace the dsn for Goose's sake with the new database name
	conf.Driver.OpenStr = strings.Replace(originalDSN, "DBNAME", dbName, -1)
	os.Setenv("DATABASE_URL", conf.Driver.OpenStr)

	// Get the migrations to run then run them
	target, err := goose.GetMostRecentDBVersion(conf.MigrationsDir)
	if err != nil {
		log.WithError(fail.Trace(err)).Fatal("Could not get most recent db version")
	}

	if err := goose.RunMigrations(conf, conf.MigrationsDir, target); err != nil {
		log.WithError(fail.Trace(err)).Fatal("Could not run setup migrations")
	}
}

// TearDown destroys the database.
func TearDown(seed string) {
	conf := newConf()

	// Delete the database entirely.
	dbName := "test_" + seed
	conf.Driver.OpenStr = strings.Replace(conf.Driver.OpenStr, "DBNAME", dbName, -1)
	db, err := sql.Open(conf.Driver.Name, conf.Driver.OpenStr)
	if err != nil {
		log.WithError(fail.Trace(err)).Fatal("Could not open db connection")
	}
	_, err = db.Query("DROP DATABASE IF EXISTS " + dbName)
	if err != nil {
		log.WithError(fail.Trace(err)).WithField("name", dbName).Fatal("Could not drop db")
	}
}

func newConf() *goose.DBConf {
	file := config.String("BASE_DIR") + config.String("DBCONF_DIR", "db")
	conf, err := goose.NewDBConf(file, "test", "")
	if err != nil {
		log.WithError(fail.Trace(err)).WithField("filename", file).Fatal("Could not create new conf")
	}
	return conf
}
