package config

import (
	"fmt"
	"os"
)

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
		Days:       1,
		IncludeAQI: true,
		Alerts:     true,
		IsLocal:    true,
	}

	cfg.APIKey = os.Getenv("WEATHER_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY environment variable is not set")
	}

	return cfg, nil
}

func (c *Config) SetLocation(location string) {
	c.Location = location
	c.IsLocal = false
}
