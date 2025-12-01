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

// Version is set at build time via -ldflags.
var version = "dev"

func main() {
	// Handle --version and -v flags
	if len(os.Args) >= 2 {
		arg := os.Args[1]
		if arg == "--version" || arg == "-v" {
			fmt.Printf("weather-cli %s\n", version)
			os.Exit(0)
		}
	}
	if err := godotenv.Load(".env"); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Allow location override via command line argument
	if len(os.Args) >= 2 {
		location := strings.TrimSpace(os.Args[1])

		if location == "" {
			fmt.Fprintf(os.Stderr, "Error: location cannot be empty\n")
			os.Exit(1)
		}

		if len(location) > 100 {
			fmt.Fprintf(os.Stderr, "Error: location too long (max 100 characters)\n")
			os.Exit(1)
		}

		cfg.SetLocation(location)
	}

	weatherCache, err := cache.New(cache.DefaultTTL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cache unavailable: %v\n", err)
	}

	var data *weather.Response
	if weatherCache != nil {
		data = weatherCache.Get(cfg.Location)
	}

	if data == nil {
		client := weather.NewClient(cfg.APIKey)
		ctx := context.Background()
		data, err = client.Fetch(ctx, weather.FetchOptions{
			Location:   cfg.Location,
			Days:       cfg.Days,
			IncludeAQI: cfg.IncludeAQI,
			Alerts:     cfg.Alerts,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching weather: %v\n", err)
			os.Exit(1)
		}

		if weatherCache != nil {
			if cacheErr := weatherCache.Set(cfg.Location, data); cacheErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to cache data: %v\n", cacheErr)
			}
		}
	}

	display, err := weatherdisplay.NewDisplay(data, cfg.IsLocal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating display: %v\n", err)
		os.Exit(1)
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
}
