package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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

	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	cfg.APIKey = os.Getenv("WEATHER_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("WEATHER_API_KEY is not set")
	}

	return cfg, nil
}

func (c *Config) SetLocation(location string) {
	c.Location = location
	c.IsLocal = false
}

func (c *Config) BuildRequestURL() string {
	urlParams := fmt.Sprintf("&days=%d", c.Days)

	if c.IncludeAQI {
		urlParams += "&aqi=yes"
	}

	if c.Alerts {
		urlParams += "&alerts=yes"
	}

	return c.BaseURL + "?key=" + c.APIKey + "&q=" + c.Location + urlParams
}
