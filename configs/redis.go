package configs

import (
	customerr "sample-web/errors"
)

type RedisConfig struct {
	Address          string `json:"address"`
	Database         int    `json:"database"`
	TimeoutInSeconds int    `json:"timeout_in_seconds"`
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
	
	return nil
}

func (redis *RedisConfig) LoadAndValidate() error {
	if err := redis.validate(); err != nil {
		return err
	}
	return nil
}
