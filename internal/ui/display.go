package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/models"
)

type Display struct {
	config  *config.Config
	weather *models.Weather
}

func NewDisplay(config *config.Config, weather *models.Weather) *Display {
	return &Display{weather: weather, config: config}
}

func (d *Display) PrintAll() {
	fmt.Print(d.Heading())
	Spacer()

	fmt.Print(d.Time())
	Spacer()

	fmt.Print(d.CurrentConditions())
	Spacer()

	fmt.Print(d.HourlyForecast())
	Spacer()

	fmt.Print(d.DailyForecast())
	Spacer()

	fmt.Print(d.Twilight())
	Spacer()

	fmt.Print(d.Warnings())
}

func (d *Display) Heading() string {
	location := d.weather.Location

	text := strings.Builder{}
	text.WriteString("\nWeather Forecast for " + location.Name + ", ")
	text.WriteString(location.Country)
	text.WriteString("\n" + CreateBorder(text.Len()))

	return text.String()
}

func (d *Display) Time() string {
	location := d.weather.Location

	localTime, err := time.Parse("2006-01-02 15:04", location.LocalTime)
	if err != nil {
		return fmt.Sprintf("Error parsing local time: %v", err)
	}

	now := time.Now()
	timeFormat := "Mon, Jan 2 - 15:04"
	timeOutput := "Time: " + now.Format(timeFormat)

	if !d.config.IsLocal {
		timeOutput += " (Local Time: " + localTime.Format(timeFormat) + ")"
	}

	return timeOutput
}

func (d *Display) CurrentConditions() string {
	c := d.weather.Current
	output := strings.Builder{}

	output.WriteString(
		"Current Conditions: " +
			GetWeatherIcon(c.Condition.Text) + " " + c.Condition.Text + ", " +
			fmt.Sprintf("%.0f°C", c.TempC) + " " +
			"(Feels like " + fmt.Sprintf("%.0f°C", c.FeelsLike) + ")\n")

	output.WriteString(
		"Wind: " +
			GetIcon("wind") + " " +
			c.WindDirection + " " +
			fmt.Sprintf("%.0f", c.WindSpeed) + " mph | ")

	output.WriteString(
		"Humidity: " +
			GetIcon("humidity") + " " +
			fmt.Sprintf("%.0f", c.Humidity) + "% | ")

	output.WriteString(
		"AQI: " +
			GetAqiIcon(c.AirQuality.PM25) + " " +
			fmt.Sprintf("%.0f", c.AirQuality.PM25) + " (PM2.5)")

	return output.String()
}

func (d *Display) HourlyForecast() string {
	output := strings.Builder{}
	output.WriteString("Hourly Forecast:\n")

	hours := d.weather.Forecast.Forecastday[0].Hour

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
				GetWeatherIcon(hour.Condition.Text),
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
	astro := d.weather.Forecast.Forecastday[0].Astro

	output := strings.Builder{}
	output.WriteString("Sunrise: " + GetIcon("sunrise") + " " + astro.Sunrise + " | ")
	output.WriteString("Sunset: " + GetIcon("sunset") + " " + astro.Sunset)

	return output.String()
}

func (d *Display) Warnings() string {
	output := strings.Builder{}
	output.WriteString("Weather Warnings: ")

	if len(d.weather.Alerts.Alert) == 0 {
		output.WriteString("None")
		return output.String()
	}

	for i, alert := range d.weather.Alerts.Alert {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(alert.Event)
	}

	return output.String()
}
