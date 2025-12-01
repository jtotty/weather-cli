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
	text.WriteString("\nWeather Forecast for ")
	text.WriteString(location.Name)
	text.WriteString(", ")
	text.WriteString(location.Country)

	// Calculate border length (excluding leading newline)
	headerLen := text.Len() - 1

	text.WriteString("\n")
	text.WriteString(ui.CreateBorder(headerLen))

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

	output.WriteString("Current Conditions: ")
	output.WriteString(ui.GetWeatherIcon(c.Condition.Text))
	output.WriteString(" ")
	output.WriteString(c.Condition.Text)
	output.WriteString(", ")
	output.WriteString(ui.ColorizeTemp(c.TempC))
	output.WriteString(" (Feels like ")
	output.WriteString(ui.ColorizeTemp(c.FeelsLike))
	output.WriteString(")\n")

	output.WriteString("Wind: ")
	output.WriteString(ui.GetIcon("wind"))
	output.WriteString(" ")
	output.WriteString(c.WindDirection)
	output.WriteString(" ")
	fmt.Fprintf(&output, "%.0f", c.WindSpeed)
	output.WriteString(" mph | ")

	output.WriteString("Humidity: ")
	output.WriteString(ui.GetIcon("humidity"))
	output.WriteString(" ")
	fmt.Fprintf(&output, "%.0f", c.Humidity)
	output.WriteString("% | ")

	output.WriteString("AQI: ")
	output.WriteString(ui.GetAqiIcon(c.AirQuality.PM25))
	output.WriteString(" ")
	fmt.Fprintf(&output, "%.0f", c.AirQuality.PM25)
	output.WriteString(" (PM2.5)")

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
	output.WriteString("Time  | Temp  | Rain | Condition\n")

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
				"%s | %s | %3.0f%% | %s %s%s",
				date.Format("15:04"),
				ui.ColorizeTemp(hour.TempC),
				hour.ChanceOfRain,
				ui.GetWeatherIcon(hour.Condition.Text),
				hour.Condition.Text,
				newLine,
			),
		)
	}

	return output.String()
}

func (d *Display) DailyForecast() string {
	if d.data == nil || len(d.data.Forecast.Forecastday) <= 1 {
		return "Daily Forecast: No data available\n"
	}

	output := strings.Builder{}
	output.WriteString("Daily Forecast:\n")
	output.WriteString("Day    | High  | Low   | Rain | Condition\n")

	// Skip today (index 0), show future days only
	forecastDays := d.data.Forecast.Forecastday[1:]
	for i := range forecastDays {
		day := &forecastDays[i]
		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			continue
		}

		output.WriteString(
			fmt.Sprintf(
				"%s | %s | %s | %3d%% | %s %s\n",
				date.Format("Mon 02"),
				ui.ColorizeTemp(day.Day.MaxTempC),
				ui.ColorizeTemp(day.Day.MinTempC),
				day.Day.ChanceOfRain,
				ui.GetWeatherIcon(day.Day.Condition.Text),
				day.Day.Condition.Text,
			),
		)
	}

	return strings.TrimSuffix(output.String(), "\n")
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
	output.WriteString("Sunrise: ")
	output.WriteString(ui.GetIcon("sunrise"))
	output.WriteString(" ")
	output.WriteString(astro.Sunrise)
	output.WriteString(" | ")
	output.WriteString("Sunset: ")
	output.WriteString(ui.GetIcon("sunset"))
	output.WriteString(" ")
	output.WriteString(astro.Sunset)

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
