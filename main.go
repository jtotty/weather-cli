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

	fmt.Print(w.Heading())
	ui.Spacer()

	fmt.Print(w.Time())
	ui.Spacer()

	fmt.Print(w.CurrentConditions())
	ui.Spacer()

	fmt.Print(w.HourlyForecast())
	ui.Spacer()

	fmt.Print(w.DailyForecast())
	ui.Spacer()

	fmt.Print(w.Twilight())
	ui.Spacer()

	fmt.Print(w.Warnings())
}
