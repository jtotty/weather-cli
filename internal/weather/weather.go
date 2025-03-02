package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jtotty/weather-cli/internal/config"
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

func Initialize() (*Weather, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	// Can pass in location as arg
	if len(os.Args) >= 2 {
		cfg.SetLocation(os.Args[1])
	}

	weather := &Weather{IsLocal: cfg.IsLocal}

	res, err := http.Get(cfg.BuildRequestURL())
	if err != nil {
		return nil, fmt.Errorf("API request to weather API failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	jsonErr := json.Unmarshal(body, &weather)
	if jsonErr != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", jsonErr)
	}

	return weather, nil
}

func (w *Weather) Heading() string {
	location := w.Location

	text := strings.Builder{}
	text.WriteString("\nWeather Forecast for " + location.Name + ", ")
	text.WriteString(location.Country)
	text.WriteString("\n" + ui.CreateBorder(text.Len()))

	return text.String()
}

func (w *Weather) Time() string {
	location := w.Location

	localTime, err := time.Parse("2006-01-02 15:04", location.LocalTime)
	if err != nil {
		return fmt.Sprintf("Error parsing local time: %v", err)
	}

	now := time.Now()
	timeFormat := "Mon, Jan 2 - 15:04"
	timeOutput := "Time: " + now.Format(timeFormat)

	if !w.IsLocal {
		timeOutput += " (Local Time: " + localTime.Format(timeFormat) + ")"
	}

	return timeOutput
}

func (w *Weather) CurrentConditions() string {
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

	return output.String()
}

func (w *Weather) HourlyForecast() string {
	output := strings.Builder{}
	output.WriteString("Hourly Forecast:\n")

	hours := w.Forecast.Forecastday[0].Hour

	// Create the table header
	output.WriteString(
		fmt.Sprintf(
			"%-5s | %-5s | %-5s | %s\n",
			"Time",
			"Temp",
			"Rain",
			"Condition",
		),
	)

	currentTime := time.Now()
	year, month, day := currentTime.Date()
	startOfNextDay := time.Date(year, month, day, 0, 0, 0, 0, currentTime.Location()).Add(24 * time.Hour)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(currentTime) {
			continue
		}

		if date.Equal(startOfNextDay) {
			break
		}

		newLine := "\n"
		if date.Add(time.Hour).Equal(startOfNextDay) {
			newLine = ""
		}

		output.WriteString(
			fmt.Sprintf(
				"%s - %.0f°C - %.0f%% - %s - %s%s",
				date.Format("15:04"),
				hour.TempC,
				hour.ChanceOfRain,
				hour.Condition.Text,
				ui.GetWeatherIcon(hour.Condition.Text),
				newLine,
			),
		)
	}

	return output.String()
}

func (w *Weather) DailyForecast() string {
	return "Daily Forecast:\n"
}

func (w *Weather) Twilight() string {
	astro := w.Forecast.Forecastday[0].Astro

	output := strings.Builder{}
	output.WriteString("Sunrise: " + ui.GetIcon("sunrise") + " " + astro.Sunrise + " | ")
	output.WriteString("Sunset: " + ui.GetIcon("sunset") + " " + astro.Sunset)

	return output.String()
}

func (w *Weather) Warnings() string {
	output := strings.Builder{}
	output.WriteString("Weather Warnings: ")

	if len(w.Alerts.Alert) == 0 {
		output.WriteString("None")
		return output.String()
	}

	for i, alert := range w.Alerts.Alert {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(alert.Event)
	}

	return output.String()
}
