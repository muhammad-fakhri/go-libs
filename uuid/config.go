package uuid

import (
	"errors"
)

type Config struct {
	baseURL  string
	clientID string
}

func NewConfig(baseURL, clientID string) (*Config, error) {
	if baseURL == "" || clientID == "" {
		return nil, errors.New("invalid config")
	}
	return &Config{baseURL, clientID}, nil
}

func (c *Config) ClientID() string {
	return c.clientID
}

func (c *Config) GetUUIClusterPath() string {
	return "/uui/cluster"
}
