package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/models"
)

type WeatherService struct {
	config *config.Config
	client *http.Client
}

func NewWeatherService(cfg *config.Config) *WeatherService {
	return &WeatherService{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *WeatherService) FetchWeather() (*models.Weather, error) {
	url := s.config.BuildRequestURL()

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var weather models.Weather
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("failed to unmarshal weather data: %w", err)
	}

	return &weather, nil
}

func (s *WeatherService) SaveToFile(weather *models.Weather, filename string) error {
	data, err := json.MarshalIndent(weather, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal weather data: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write weather data to file: %w", err)
	}

	return nil
}
