package main

import (
	"github.com/jtotty/weather-cli/internal/ui"
	"github.com/jtotty/weather-cli/internal/weather"
)

func main() {
	w := weather.Initialize()
	w.Heading()
	ui.Spacer()

	w.CurrentConditions()
	ui.Spacer()

	w.HourlyForecast()
	ui.Spacer()

	w.Warnings()
}
