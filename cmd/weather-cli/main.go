package main

import (
	"github.com/jtotty/weather-cli/internal/ui"
	"github.com/jtotty/weather-cli/internal/weather"
)

func main() {
    w := weather.GenerateWeather()
	now := w.Now()
	ui.SingleFrameDisplay(now, len(now))

	hours, hoursMaxLen := w.Hours()
	ui.MultilineFrameDisplay(hours, hoursMaxLen)
}
