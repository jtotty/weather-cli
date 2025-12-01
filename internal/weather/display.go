package weather

import (
	"fmt"
	"strings"
	"time"

	api "github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/ui"
)

type Display struct {
	data    *api.Response
	isLocal bool
}

func NewDisplay(data *api.Response, isLocal bool) (*Display, error) {
	if data == nil {
		return nil, fmt.Errorf("weather data is nil")
	}

	if len(data.Forecast.Forecastday) == 0 {
		return nil, fmt.Errorf("no forecast data available")
	}

	return &Display{
		data:    data,
		isLocal: isLocal,
	}, nil
}

func (d *Display) Heading() string {
	location := d.data.Location

	text := strings.Builder{}
	text.WriteString("\nWeather Forecast for " + location.Name + ", ")
	text.WriteString(location.Country)
	text.WriteString("\n" + ui.CreateBorder(text.Len()))

	return text.String()
}

func (d *Display) Time() string {
	if d.data == nil || d.data.Location == (api.Location{}) {
		return "Time: No data available\n"
	}

	location := d.data.Location
	timeFormat := "Mon, Jan 2 - 15:04"
	now := time.Now()

	localTime, err := time.Parse("2006-01-02 15:04", location.LocalTime)
	if err != nil {
		return "Time: " + now.Format(timeFormat)
	}

	timeOutput := "Time: " + now.Format(timeFormat)

	if !d.isLocal {
		timeOutput += " (Local Time: " + localTime.Format(timeFormat) + ")"
	}

	return timeOutput
}

func (d *Display) CurrentConditions() string {
	if d.data == nil || d.data.Current == (api.Current{}) {
		return "Current Conditions: No data available\n"
	}

	c := d.data.Current
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

func (d *Display) HourlyForecast() string {
	if d.data == nil || len(d.data.Forecast.Forecastday) == 0 {
		return "Hourly Forecast: No data available\n"
	}

	hours := d.data.Forecast.Forecastday[0].Hour
	if len(hours) == 0 {
		return "Hourly Forecast: No hourly data available\n"
	}

	output := strings.Builder{}
	output.WriteString("Hourly Forecast:\n")

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

func (d *Display) DailyForecast() string {
	return "Daily Forecast:\n"
}

func (d *Display) Twilight() string {
	if d.data == nil || len(d.data.Forecast.Forecastday) == 0 {
		return "Twilight: No data available\n"
	}

	astro := d.data.Forecast.Forecastday[0].Astro
	if astro.Sunrise == "" || astro.Sunset == "" {
		return "Twilight: No sunrise or sunset data available\n"
	}

	output := strings.Builder{}
	output.WriteString("Sunrise: " + ui.GetIcon("sunrise") + " " + astro.Sunrise + " | ")
	output.WriteString("Sunset: " + ui.GetIcon("sunset") + " " + astro.Sunset)

	return output.String()
}

func (d *Display) Warnings() string {
	if d.data == nil || len(d.data.Alerts.Alert) == 0 {
		return "Weather Warnings: None\n"
	}

	output := strings.Builder{}
	output.WriteString("Weather Warnings: ")

	for i, alert := range d.data.Alerts.Alert {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(alert.Event)
	}

	return output.String()
}
