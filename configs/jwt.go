package configs

import (
	"os"
	customerr "sample-web/errors"
)

type JWTConfig struct {
	ExpirationInSeconds             int    `json:"expiration_in_seconds"`
	IssuerName                      string `json:"issuer_name"`
	SecretKey                       string `json:"-"`
	RefreshTokenExpirationInSeconds int    `json:"refresh_token_expiration_in_seconds"`
	RefreshTokenSecret              string `json:"-"`
}

func (jwtConfig *JWTConfig) validate() error {
	if jwtConfig.ExpirationInSeconds <= 0 {
		return customerr.MissingConfigError{Message: "expiration_in_seconds must be greater than 0"}
	}
	if jwtConfig.IssuerName == "" {
		return customerr.MissingConfigError{Message: "issuer_name is not set"}
	}
	if jwtConfig.SecretKey == "" {
		return customerr.MissingConfigError{Message: "JWT_SECRET_KEY is not set"}
	}
	if jwtConfig.RefreshTokenExpirationInSeconds <= 0 {
		return customerr.MissingConfigError{Message: "refresh_token_expiration_in_seconds must be greater than 0"}
	}
	if jwtConfig.RefreshTokenSecret == "" {
		return customerr.MissingConfigError{Message: "JWT_REFRESH_SECRET_KEY is not set"}
	}
	return nil
}

func (jwtConfig *JWTConfig) LoadAndValidate() error {
	jwtConfig.SecretKey = os.Getenv("JWT_SECRET_KEY")
	jwtConfig.RefreshTokenSecret = os.Getenv("JWT_REFRESH_SECRET_KEY")
	if err := jwtConfig.validate(); err != nil {
		return err
	}
	return nil
}
