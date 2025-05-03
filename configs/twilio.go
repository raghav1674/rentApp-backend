package configs

import (
	"os"
	customerr "sample-web/errors"
)

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	ServiceSID string
}

func (config *TwilioConfig) validate() error {
	if config.AccountSID == "" {
		return customerr.MissingConfigError{Message: "TWILIO_ACCOUNT_SID is not set"}
	}
	if config.AuthToken == "" {
		return customerr.MissingConfigError{Message: "TWILIO_AUTH_TOKEN is not set"}
	}
	if config.ServiceSID == "" {
		return customerr.MissingConfigError{Message: "TWILIO_VERIFY_SERVICE_SID is not set"}
	}
	return nil
}

func (config *TwilioConfig) LoadAndValidate() error {
	config.AccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	config.AuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	config.ServiceSID = os.Getenv("TWILIO_VERIFY_SERVICE_SID")
	if err := config.validate(); err != nil {
		return err
	}
	return nil
}
