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
		Name      string `json:"name"`
		Country   string `json:"country"`
		LocalTime string `json:"localtime"`
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
			AirQuality struct {
				PM25 float32 `json:"pm2_5"`
				PM10 float32 `json:"pm10"`
			} `json:"air_quality"`
			Astro struct {
				Sunrise string `json:"sunrise"`
				Sunset  string `json:"sunset"`
			} `json:"astro"`
		} `json:"forecastday"`
	} `json:"forecast"`
	Alerts struct {
		Alert []struct {
			Event string `json:"event"`
			Desc  string `json:"desc"`
		} `json:"alert"`
	} `json:"alerts"`
	IsLocal bool
}

func Initialize() *Weather {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		fmt.Println("ENVIRONMENT ERROR: WEATHER_API_KEY is not set")
		os.Exit(1)
	}

	baseURL := "https://api.weatherapi.com/v1/forecast.json"
	location := "auto:ip"
	options := "&days=1&aqi=yes&alerts=yes"

	var weather Weather
	weather.IsLocal = true

	// Can pass in location as arg
	if len(os.Args) >= 2 {
		location = os.Args[1]
		weather.IsLocal = false
	}

	res, err := http.Get(baseURL + "?key=" + apiKey + "&q=" + location + options)
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

	jsonErr := json.Unmarshal(body, &weather)
	if jsonErr != nil {
		panic(jsonErr)
	}

	return &weather
}

func (w *Weather) Heading() {
	location := w.Location

	text := strings.Builder{}
	text.WriteString("\nWeather Forecast for " + location.Name + ", ")
	text.WriteString(location.Country)
	text.WriteString("\n" + ui.CreateBorder(text.Len()))

	fmt.Print(text.String())
}

func (w *Weather) Time() {
	location := w.Location

	localTime, err := time.Parse("2006-01-02 15:04", location.LocalTime)
	if err != nil {
		panic("Error parsing local time")
	}

	now := time.Now()
	timeFormat := "Mon, Jan 2 - 15:04"
	timeOutput := "Time: " + now.Format(timeFormat)

	if !w.IsLocal {
		timeOutput += " (Local Time: " + localTime.Format(timeFormat) + ")"
	}

	fmt.Print(timeOutput)
}

func (w *Weather) CurrentConditions() {
	c := w.Current
	output := strings.Builder{}

	output.WriteString(
		"Current Conditions: " +
			ui.GetWeatherIcon(c.Condition.Text) + " " + c.Condition.Text + ", " +
			fmt.Sprintf("%.0f°C", c.TempC) + " " +
			"(Feels like " + fmt.Sprintf("%.0f°C", c.FeelsLike) + ")\n")

	output.WriteString(
		"Wind: " +
			ui.GetIcon("wind") + " " +
			c.WindDirection + " " +
			fmt.Sprintf("%.0f", c.WindSpeed) + " mph | ")

	output.WriteString(
		"Humidity: " +
			ui.GetIcon("humidity") + " " +
			fmt.Sprintf("%.0f", c.Humidity) + "% | ")

	output.WriteString(
		"AQI: " +
			ui.GetAqiIcon(c.AirQuality.PM25) + " " +
			fmt.Sprintf("%.0f", c.AirQuality.PM25) + " (PM2.5)")

	fmt.Print(output.String())
}

func (w *Weather) HourlyForecast() {
	fmt.Println("Hourly Forecast:")

	hours := w.Forecast.Forecastday[0].Hour

	// Create the table header
	fmt.Printf(
		"%-5s | %-5s | %-5s | %s\n",
		"Time",
		"Temp",
		"Rain",
		"Condition",
	)

	currentTime := time.Now()
	year, month, day := currentTime.Date()
	startOfNextDay := time.Date(year, month, day, 0, 0, 0, 0, currentTime.Location()).Add(24 * time.Hour)

	newLine := "\n"

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(currentTime) {
			continue
		}

		if date.Equal(startOfNextDay) {
			break
		}

		if date.Add(time.Hour).Equal(startOfNextDay) {
			newLine = ""
		}

		fmt.Printf(
			"%s - %.0f°C - %.0f%% - %s - %s"+newLine,
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
			ui.GetWeatherIcon(hour.Condition.Text),
		)
	}
}

func (w *Weather) DailyForecast() {
	fmt.Println("Daily Forecast:")
}

func (w *Weather) Twilight() {
	astro := w.Forecast.Forecastday[0].Astro

	output := strings.Builder{}
	output.WriteString("Sunrise: " + ui.GetIcon("sunrise") + " " + astro.Sunrise + " | ")
	output.WriteString("Sunset: " + ui.GetIcon("sunset") + " " + astro.Sunset)

	fmt.Print(output.String())
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
