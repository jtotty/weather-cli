package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jtotty/weather-cli/internal/ui"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC         float32 `json:"temp_c"`
		FeelsLike     float32 `json:"feelslike_c"`
		Humidity      float32 `json:"humidity"`
		WindSpeed     float32 `json:"wind_mph"`
		WindDirection string  `json:"wind_dir"`
		Condition     struct {
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
	Alerts struct {
		Alert []struct {
			Event string `json:"event"`
			Desc  string `json:"desc"`
		} `json:"alert"`
	} `json:"alerts"`
}

func Initialize() *Weather {
	baseURL := "https://api.weatherapi.com/v1/forecast.json"
	key := "7937cf0616e0430aaf534238241701"
	location := "auto:ip"
	options := "&days=1&aqi=yes&alerts=yes"

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

	var weather Weather
	jsonErr := json.Unmarshal(body, &weather)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return &weather
}

func (w *Weather) Heading() {
	location := w.Location
	now := time.Now()

	text := strings.Builder{}
	text.WriteString(location.Name + ", ")
	text.WriteString(location.Country + " | " + now.Format("Mon, Jan 2 - 15:04"))

	prepend := "\nWeather Forecast for "
	border := ui.CreateBorder(text.Len() + len(prepend))

	formatted := prepend +
		text.String() + "\n" +
		border

	fmt.Println(formatted)
}

func (w *Weather) CurrentConditions() {
	c := w.Current

	output := strings.Builder{}
	output.WriteString("Current Conditions: " + c.Condition.Text + ", ")
	output.WriteString(fmt.Sprintf("%.0f°C", c.TempC) + " ")
	output.WriteString("(Feels like " + fmt.Sprintf("%.0f°C", c.FeelsLike) + ")\n")
	output.WriteString("Wind: " + c.WindDirection + " " + fmt.Sprintf("%.0f", c.WindSpeed) + " mph | ")
	output.WriteString("Humidity: " + fmt.Sprintf("%.0f", c.Humidity) + "% | ")
	output.WriteString("Polution: " + "pm2.5 " + fmt.Sprintf("%.0f", c.AirQuality.PM25) + " ")
	output.WriteString("pm10 " + fmt.Sprintf("%.0f", c.AirQuality.PM10) + "\n")

	fmt.Print(output.String())
}

// Hourly weather data after the current time up to 23:00 hours
func (w *Weather) HourlyForecast() {
	fmt.Println("Houry Forecast:")

	hours := w.Forecast.Forecastday[0].Hour

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		// Only display even hours in the future
		if date.Before(time.Now()) {
			continue
		}

		fmt.Printf(
			"%s - %.0f°C - %s - %.0f%%\n",
			date.Format("15:04"),
			hour.TempC,
			hour.Condition.Text,
			hour.ChanceOfRain,
		)
	}
}

func (w *Weather) DailyForecast() {

}

func (w *Weather) Sun() {

}

func (w *Weather) Warnings() {
	fmt.Print("Weather Warnings: ")

	if len(w.Alerts.Alert) == 0 {
		fmt.Println("None")
		return
	}

	for _, alert := range w.Alerts.Alert {
		fmt.Println(alert.Event)
	}
}
