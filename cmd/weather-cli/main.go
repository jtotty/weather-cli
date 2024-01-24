package main

import "github.com/jtotty/weather-cli/internal/weather"

func main() {
	json := weather.QueryAPI()
    data := weather.CreateWeather(json)

    weather.Display(data)
}
