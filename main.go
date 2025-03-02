package main

import (
	"fmt"
	"os"

	"github.com/jtotty/weather-cli/internal/ui"
	"github.com/jtotty/weather-cli/internal/weather"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(1)
	}

	w, err := weather.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing weather service: %v\n", err)
		os.Exit(1)
	}

	w.Heading()
	ui.Spacer()

	w.Time()
	ui.Spacer()

	w.CurrentConditions()
	ui.Spacer()

	w.HourlyForecast()
	ui.Spacer()

	w.DailyForecast()
	ui.Spacer()

	w.Twilight()
	ui.Spacer()

	w.Warnings()
}
