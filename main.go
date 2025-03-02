package main

import (
	"fmt"
	"os"

	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/services"
	"github.com/jtotty/weather-cli/internal/ui"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	weatherService := services.NewWeatherService(cfg)

	weather, err := weatherService.FetchWeather()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing weather service: %v\n", err)
		os.Exit(1)
	}

	display := ui.NewDisplay(cfg, weather)
	display.PrintAll()
}
