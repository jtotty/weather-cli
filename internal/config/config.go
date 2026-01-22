package config

import (
	"github.com/jtotty/weather-cli/internal/credentials"
)

// Config holds the application configuration.
type Config struct {
	APIKey     string
	Location   string
	Days       int
	IncludeAQI bool
	Alerts     bool
	IsLocal    bool
}

// NewWithAPIKey creates a Config with the provided API key.
// Use this for testing or when the API key is obtained externally.
func NewWithAPIKey(apiKey string) *Config {
	return &Config{
		APIKey:     apiKey,
		Location:   "auto:ip",
		Days:       7,
		IncludeAQI: true,
		Alerts:     true,
		IsLocal:    true,
	}
}

// New creates a Config by loading the API key from credentials.
func New() (*Config, error) {
	apiKey, err := credentials.GetAPIKey()
	if err != nil {
		return nil, err
	}
	return NewWithAPIKey(apiKey), nil
}

func (c *Config) SetLocation(location string) {
	c.Location = location
	c.IsLocal = false
}
