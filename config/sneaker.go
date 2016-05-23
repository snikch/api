package config

import (
	"bytes"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/codahale/sneaker"
	"github.com/snikch/api/log"
)

// initSneakerSecrets loads any remote secrets from S3.
func initSneakerSecrets() {
	secrets, err := RemoteSecrets()
	if err != nil {
		log.WithError(err).Fatal("Could not download secrets")
	}
	log.WithField("total", len(secrets)).Info("Found remote secrets")
	for name, value := range secrets {
		log.WithField("key", name).Info("Setting secret")
		os.Setenv(name, string(value))
	}
}

// RemoteSecrets returns all remove secrets as key value pairs.
func RemoteSecrets() (map[string][]byte, error) {
	manager := loadSneakerManager()
	if manager == nil {
		return nil, nil
	}
	// If we have a sneaker manager, we should list it and download all creds.
	files, err := manager.List("")
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = file.Path
	}
	return manager.Download(paths)
}

// SetRemoteSecret sets a remote config value.
func SetRemoteSecret(key, value string) error {
	manager := loadSneakerManager()
	if manager == nil {
		return fmt.Errorf("Unable to create sneaker manager")
	}
	return manager.Upload(key, bytes.NewBuffer([]byte(value)))
}

// loadSneakerManager returns a sneaker manager if one is configured on the env.
func loadSneakerManager() *sneaker.Manager {
	path := String("SNEAKER_S3_PATH")
	if path == "" {
		return nil
	}
	u, err := url.Parse(path)
	if err != nil {
		log.WithField("path", path).Fatal("Invalid SNEAKER_S3_PATH")
		return nil
	}
	if u.Path != "" && u.Path[0] == '/' {
		u.Path = u.Path[1:]
	}

	sess := session.New()
	config := aws.NewConfig().WithRegion(String("SNEAKER_S3_REGION", "us-west-2")).WithMaxRetries(3)
	return &sneaker.Manager{
		Objects: s3.New(sess, config),
		Envelope: sneaker.Envelope{
			KMS: kms.New(sess, config),
		},
		Bucket: u.Host,
		Prefix: u.Path,
		KeyId:  String("SNEAKER_MASTER_KEY"),
	}
}
