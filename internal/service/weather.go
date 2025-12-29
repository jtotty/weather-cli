// Package service provides the weather service that orchestrates fetching and caching.
package service

import (
	"context"
	"fmt"
	"os"

	"github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/cache"
	"github.com/jtotty/weather-cli/internal/config"
)

// WeatherFetcher defines the interface for fetching weather data.
type WeatherFetcher interface {
	Fetch(ctx context.Context, opts weather.FetchOptions) (*weather.Response, error)
}

// WeatherCache defines the interface for caching weather data.
type WeatherCache interface {
	Get(location string) *weather.Response
	Set(location string, data *weather.Response) error
}

// Weather orchestrates fetching weather data with caching.
type Weather struct {
	cfg     *config.Config
	cache   WeatherCache
	fetcher WeatherFetcher
}

// NewWeather creates a new Weather service with default cache and API client.
func NewWeather(cfg *config.Config) *Weather {
	weatherCache, err := cache.New(cache.DefaultTTL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cache unavailable: %v\n", err)
	}

	var cacheImpl WeatherCache
	if weatherCache != nil {
		cacheImpl = weatherCache
	}

	return &Weather{
		cfg:     cfg,
		cache:   cacheImpl,
		fetcher: weather.NewClient(cfg.APIKey),
	}
}

// NewWeatherWithDeps creates a Weather service with injected dependencies (for testing).
func NewWeatherWithDeps(cfg *config.Config, c WeatherCache, fetcher WeatherFetcher) *Weather {
	return &Weather{
		cfg:     cfg,
		cache:   c,
		fetcher: fetcher,
	}
}

// GetWeather retrieves weather data, using cache if available.
func (w *Weather) GetWeather() (*weather.Response, error) {
	if w.cache != nil {
		if data := w.cache.Get(w.cfg.Location); data != nil {
			return data, nil
		}
	}

	data, err := w.fetchFromAPI()
	if err != nil {
		return nil, err
	}

	if w.cache != nil {
		if cacheErr := w.cache.Set(w.cfg.Location, data); cacheErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to cache data: %v\n", cacheErr)
		}
	}

	return data, nil
}

func (w *Weather) fetchFromAPI() (*weather.Response, error) {
	ctx := context.Background()

	return w.fetcher.Fetch(ctx, weather.FetchOptions{
		Location:   w.cfg.Location,
		Days:       w.cfg.Days,
		IncludeAQI: w.cfg.IncludeAQI,
		Alerts:     w.cfg.Alerts,
	})
}
