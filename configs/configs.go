package configs

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Mongo MongoConfig `json:"mongo"`
	Env   string      `json:"env"`
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

	return cfg, nil
}

// GetMongoConfig returns the Mongo configuration
func (config *Config) GetMongoConfig() MongoConfig {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.Mongo
}

// GetEnv returns the value of the environment
func (config *Config) GetEnvironment() string {
	if config == nil {
		panic("Config not loaded. Call LoadConfig() first.")
	}
	return config.Env
}
