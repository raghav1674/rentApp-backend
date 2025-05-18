package configs

import (
	"os"
	customerr "sample-web/errors"
)

type RedisConfig struct {
	Address          string `json:"address"`
	Database         int    `json:"database"`
	TimeoutInSeconds int    `json:"timeout_in_seconds"`
	Username         string `json:"username"`
	Password		 string `json:"password"`
	AuthEnabled	     bool   `json:"auth_enabled"`
}

func (redis *RedisConfig) validate() error {
	if redis.Address == "" {
		return customerr.MissingConfigError{Message: "redis address is not set"}
	}

	if redis.Database > 16 || redis.Database < 0 {
		return customerr.MissingConfigError{Message: "redis database must be between 0 and 15"}
	}

	if redis.TimeoutInSeconds <= 0 {
		return customerr.MissingConfigError{Message: "redis timeout_in_seconds must be greater than 0"}
	}
	if redis.AuthEnabled {
		if redis.Username == "" {
			return customerr.MissingConfigError{Message: "REDIS_USERNAME is not set"}
		}
		if redis.Password == "" {
			return customerr.MissingConfigError{Message: "REDIS_PASSWORD is not set"}
		}
	}
	return nil
}

func (redis *RedisConfig) LoadAndValidate() error {
	redis.Username = os.Getenv("REDIS_USERNAME")
	redis.Password = os.Getenv("REDIS_PASSWORD")
	if err := redis.validate(); err != nil {
		return err
	}
	return nil
}
