package config

import (
	"fmt"
	"os"
)

type Config struct {
	APIKey     string
	BaseURL    string
	Location   string
	Days       int
	IncludeAQI bool
	Alerts     bool
	IsLocal    bool
}

func New() (*Config, error) {
	cfg := &Config{
		BaseURL:    "https://api.weatherapi.com/v1/forecast.json",
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

func (c *Config) BuildRequestURL() string {
	options := fmt.Sprintf("&days=%d", c.Days)

	if c.IncludeAQI {
		options += "&aqi=yes"
	}

	if c.Alerts {
		options += "&alerts=yes"
	}

	return c.BaseURL + "?key=" + c.APIKey + "&q=" + c.Location + options
}
