package configs

import (
	"os"
	customerr "sample-web/errors"
)

type MongoConfig struct {
	URI              string `json:"uri"`
	Database         string `json:"database"`
	TimeoutInSeconds int    `json:"timeout_in_seconds"`
	Username         string `json:"-"` // will be populated from env
	Password         string `json:"-"` // will be populated from env
	AuthSource       string `json:"auth_source"`
}

func (mongo *MongoConfig) validate() error {
	if mongo.URI == "" {
		return customerr.MissingConfigError{Message: "MONGO_URI is not set"}
	}
	if mongo.Database == "" {
		return customerr.MissingConfigError{Message: "MONGO_DATABASE is not set"}
	}
	if mongo.TimeoutInSeconds <= 0 {
		return customerr.MissingConfigError{Message: "MONGO_TIMEOUT_IN_SECONDS must be greater than 0"}
	}
	if mongo.Username == "" || mongo.Password == "" {
		return customerr.MissingConfigError{Message: "MONGO_USERNAME or MONGO_PASSWORD is not set"}
	}
	return nil
}

func (mongo *MongoConfig) LoadAndValidate() error {
	mongo.Username = os.Getenv("MONGO_USERNAME")
	mongo.Password = os.Getenv("MONGO_PASSWORD")
	if err := mongo.validate(); err != nil {
		return err
	}
	return nil
}
