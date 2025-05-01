package configs

import (
	customerr "sample-web/errors"
)

type TracingConfig struct {
	ServiceName  string `json:"service_name"`
	CollectorUrl string `json:"collector_url"`
	Insecure     bool   `json:"insecure"`
}

func (tracingConfig *TracingConfig) validate() error {
	if tracingConfig.ServiceName == "" {
		return customerr.MissingConfigError{Message: "service_name must not be empty"}
	}
	if tracingConfig.CollectorUrl == "" {
		return customerr.MissingConfigError{Message: "collector_url must not be empty"}
	}
	return nil
}

func (tracingConfig *TracingConfig) LoadAndValidate() error {
	if err := tracingConfig.validate(); err != nil {
		return err
	}
	return nil
}
