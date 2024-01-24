package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

type Data struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float32 `json:"temp_c"`
		FeelsLike float32 `json:"feelslike_c"`
		Humidity  float32 `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		AirQuality struct {
			PM25 float32 `json:"pm2_5"`
			PM10 float32 `json:"pm10"`
		} `json:"air_quality"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float32 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float32 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func QueryAPI() []byte {
	baseURL := "http://api.weatherapi.com/v1/forecast.json"
	key := "7937cf0616e0430aaf534238241701"
	location := "auto:ip"
	options := "&days=1&aqi=yes&alerts=no"

	// Can pass in location as arg
	if len(os.Args) >= 2 {
		location = os.Args[1]
	}

	res, err := http.Get(baseURL + "?key=" + key + "&q=" + location + options)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available...")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

    return body
}

func CreateWeather(body []byte) Data {
	var weather Data

    err := json.Unmarshal(body, &weather)
    if err != nil {
        panic(err)
    }

	return weather
}

func Display(weather Data) {
	location := weather.Location
	current := weather.Current
	hours := weather.Forecast.Forecastday[0].Hour

	fmt.Printf(
		"%s, %s: %.0fC, %s, PM2.5 %.1f, PM10 %.1f\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
		current.AirQuality.PM25,
		current.AirQuality.PM10,
	)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		// Only display hours in the future
		if date.Before(time.Now()) {
			continue
		}

		message := fmt.Sprintf(
			"%s -- %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Print(message)
		} else {
			color.Blue(message)
		}
	}
}
