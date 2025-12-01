package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/ui"
	weatherdisplay "github.com/jtotty/weather-cli/internal/weather"

	"github.com/joho/godotenv"
)

func main() {
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

	client := weather.NewClient(cfg.APIKey)
	ctx := context.Background()
	data, err := client.Fetch(ctx, weather.FetchOptions{
		Location:   cfg.Location,
		Days:       cfg.Days,
		IncludeAQI: cfg.IncludeAQI,
		Alerts:     cfg.Alerts,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching weather: %v\n", err)
		os.Exit(1)
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
