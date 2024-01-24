package main

import (
	"github.com/jtotty/weather-cli/internal/ui"
	"github.com/jtotty/weather-cli/internal/weather"
)

func main() {
	json := weather.QueryAPI()
    data := weather.CreateWeather(json)

    current := weather.Current(data)
    ui.FrameDisplay(current)
}
