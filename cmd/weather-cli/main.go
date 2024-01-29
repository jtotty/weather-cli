package main

import (
	"github.com/jtotty/weather-cli/internal/ui"
	"github.com/jtotty/weather-cli/internal/weather"
)

func main() {
	json := weather.QueryAPI()
	data := weather.CreateWeather(json)

	now := weather.Now(&data)
	ui.SingleFrameDisplay(now, len(now))

	// hours, hoursMaxLen := weather.Hours(&data)
	// ui.MultilineFrameDisplay(hours, hoursMaxLen)
}
