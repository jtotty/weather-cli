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
		fmt.Println("ENVIRONMENT ERROR: Error loading .env file")
		os.Exit(1)
	}

	w := weather.Initialize()

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
