package main

import "github.com/jtotty/weather-cli/internal/weather"

func main() {
	data := weather.QueryAPI()
    weather.Display(data)
}
