package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/cache"
	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/ui"
	weatherdisplay "github.com/jtotty/weather-cli/internal/weather"

	"github.com/joho/godotenv"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	if handleVersionFlag() {
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		exitWithErrorf("Error loading .env file: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		exitWithErrorf("Error loading config: %v", err)
	}

	if err := parseLocationArg(cfg); err != nil {
		exitWithErrorf("%v", err)
	}

	data, err := fetchWeatherData(cfg)
	if err != nil {
		exitWithErrorf("Error fetching weather: %v", err)
	}

	if err := displayWeather(data, cfg.IsLocal); err != nil {
		exitWithErrorf("Error creating display: %v", err)
	}
}

// handleVersionFlag checks for --version or -v flags and prints version if found.
func handleVersionFlag() bool {
	if len(os.Args) >= 2 {
		arg := os.Args[1]
		if arg == "--version" || arg == "-v" {
			fmt.Printf("weather-cli %s\n", version)
			return true
		}
	}
	return false
}

// parseLocationArg parses the location from command line arguments.
func parseLocationArg(cfg *config.Config) error {
	if len(os.Args) < 2 {
		return nil
	}

	location := strings.TrimSpace(os.Args[1])

	if location == "" {
		return fmt.Errorf("location cannot be empty")
	}

	if len(location) > 100 {
		return fmt.Errorf("location too long (max 100 characters)")
	}

	cfg.SetLocation(location)
	return nil
}

// fetchWeatherData retrieves weather data from cache or API.
func fetchWeatherData(cfg *config.Config) (*weather.Response, error) {
	weatherCache, err := cache.New(cache.DefaultTTL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cache unavailable: %v\n", err)
	}

	// Try cache first
	if weatherCache != nil {
		if data := weatherCache.Get(cfg.Location); data != nil {
			return data, nil
		}
	}

	// Fetch from API
	data, err := fetchFromAPI(cfg)
	if err != nil {
		return nil, err
	}

	// Cache the response
	if weatherCache != nil {
		if cacheErr := weatherCache.Set(cfg.Location, data); cacheErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to cache data: %v\n", cacheErr)
		}
	}

	return data, nil
}

// fetchFromAPI fetches weather data from the weather API.
func fetchFromAPI(cfg *config.Config) (*weather.Response, error) {
	client := weather.NewClient(cfg.APIKey)
	ctx := context.Background()

	return client.Fetch(ctx, weather.FetchOptions{
		Location:   cfg.Location,
		Days:       cfg.Days,
		IncludeAQI: cfg.IncludeAQI,
		Alerts:     cfg.Alerts,
	})
}

// displayWeather renders the weather data to stdout.
func displayWeather(data *weather.Response, isLocal bool) error {
	display, err := weatherdisplay.NewDisplay(data, isLocal)
	if err != nil {
		return err
	}

	fmt.Print(display.Heading())
	ui.Spacer()

	fmt.Print(display.Time())
	ui.Spacer()

	fmt.Print(display.CurrentConditions())
	ui.Spacer()

	fmt.Print(display.HourlyForecast())
	ui.Spacer()

	fmt.Print(display.DailyForecast())
	ui.Spacer()

	fmt.Print(display.Twilight())
	ui.Spacer()

	fmt.Print(display.Warnings())

	return nil
}

// exitWithErrorf prints an error message to stderr and exits with code 1.
func exitWithErrorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
