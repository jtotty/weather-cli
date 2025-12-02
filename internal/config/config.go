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

func New() (*Config, error) {
	cfg := &Config{
		Location:   "auto:ip",
		Days:       7,
		IncludeAQI: true,
		Alerts:     true,
		IsLocal:    true,
	}

	apiKey, err := credentials.GetAPIKey()
	if err != nil {
		return nil, err
	}

	cfg.APIKey = apiKey
	return cfg, nil
}

func (c *Config) SetLocation(location string) {
	c.Location = location
	c.IsLocal = false
}
