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

type Weather struct {
	cfg   *config.Config
	cache *cache.Cache
}

func NewWeather(cfg *config.Config) *Weather {
	weatherCache, err := cache.New(cache.DefaultTTL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cache unavailable: %v\n", err)
	}

	return &Weather{
		cfg:   cfg,
		cache: weatherCache,
	}
}

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
	client := weather.NewClient(w.cfg.APIKey)
	ctx := context.Background()

	return client.Fetch(ctx, weather.FetchOptions{
		Location:   w.cfg.Location,
		Days:       w.cfg.Days,
		IncludeAQI: w.cfg.IncludeAQI,
		Alerts:     w.cfg.Alerts,
	})
}
