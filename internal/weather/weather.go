package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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

func Now(weather Data) string {
	location := weather.Location
	current := weather.Current

	text := strings.Builder{}

	text.WriteString(location.Name + ", ")
	text.WriteString(location.Country + ": ")
	text.WriteString(fmt.Sprintf("%.0fC", current.TempC) + ", ")
	text.WriteString(current.Condition.Text + ", ")
	text.WriteString("PM2.5 " + fmt.Sprintf("%.1f", current.AirQuality.PM25) + ", ")
	text.WriteString("PM10 " + fmt.Sprintf("%.1f", current.AirQuality.PM10))

	return text.String()
}

func Hours(weather Data) (string, int) {
	rows := []string{}
	hours := weather.Forecast.Forecastday[0].Hour
	longestStr := 0

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		// Only display hours in the future
		if date.Before(time.Now()) {
			continue
		}

		formatted := fmt.Sprintf(
			"%s -- %.0fC, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		len := len(formatted)
		if len > longestStr {
			longestStr = len
		}

		rows = append(rows, formatted)
	}

	return strings.Join(rows, ""), longestStr
}
