package configs

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Mongo   MongoConfig   `json:"mongo"`
	JWT     JWTConfig     `json:"jwt"`
	Env     string        `json:"env"`
	Tracing TracingConfig `json:"tracing"`
	Twilio  TwilioConfig  `json:"twilio"`
	Redis   RedisConfig   `json:"redis"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if err := cfg.Mongo.LoadAndValidate(); err != nil {
		return nil, err
	}

	if err := cfg.JWT.LoadAndValidate(); err != nil {
		return nil, err
	}

	if err := cfg.Tracing.LoadAndValidate(); err != nil {
		return nil, err
	}

	if err := cfg.Twilio.LoadAndValidate(); err != nil {
		return nil, err
	}

	if err := cfg.Redis.LoadAndValidate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetMongoConfig returns the Mongo configuration
func (config *Config) GetMongoConfig() MongoConfig {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.Mongo
}

// GetJWTConfig returns the JWT configuration
func (config *Config) GetJWTConfig() JWTConfig {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.JWT
}

// GetTracingConfig returns the Tracing configuration
func (config *Config) GetTracingConfig() TracingConfig {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.Tracing
}

// GetTwilioConfig returns the Twilio configuration
func (config *Config) GetTwilioConfig() TwilioConfig {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.Twilio
}

// GetRedisConfig returns the Redis configuration
func (config *Config) GetRedisConfig() RedisConfig {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.Redis
}

// GetEnv returns the value of the environment
func (config *Config) GetEnvironment() string {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.Env
}
