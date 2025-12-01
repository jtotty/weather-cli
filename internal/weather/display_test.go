package weather

import (
	"strings"
	"testing"

	api "github.com/jtotty/weather-cli/internal/api/weather"
)

func TestNewDisplay_Validation(t *testing.T) {
	tests := []struct {
		name    string
		data    *api.Response
		wantErr string
	}{
		{
			name:    "nil data",
			data:    nil,
			wantErr: "nil",
		},
		{
			name: "empty forecast",
			data: &api.Response{
				Forecast: api.Forecast{Forecastday: []api.ForecastDay{}},
			},
			wantErr: "no forecast",
		},
		{
			name: "valid data",
			data: &api.Response{
				Forecast: api.Forecast{Forecastday: []api.ForecastDay{{}}},
			},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			display, err := NewDisplay(tt.data, true)

			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if display == nil {
					t.Error("expected display, got nil")
				}
			} else {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if err != nil && !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %q, want error containing %q", err.Error(), tt.wantErr)
				}
				if display != nil {
					t.Error("expected nil display when error occurs")
				}
			}
		})
	}
}

func TestHourlyForecast_EmptyHours(t *testing.T) {
	data := &api.Response{
		Forecast: api.Forecast{
			Forecastday: []api.ForecastDay{
				{Hour: []api.Hour{}},
			},
		},
	}

	display, err := NewDisplay(data, true)
	if err != nil {
		t.Fatalf("unexpected error creating display: %v", err)
	}

	result := display.HourlyForecast()

	if !strings.Contains(result, "No hourly data") {
		t.Errorf("HourlyForecast() = %q, want string containing 'No hourly data'", result)
	}
}

func TestTwilight_MissingData(t *testing.T) {
	data := &api.Response{
		Forecast: api.Forecast{
			Forecastday: []api.ForecastDay{
				{Astro: api.Astro{Sunrise: "", Sunset: ""}},
			},
		},
	}

	display, err := NewDisplay(data, true)
	if err != nil {
		t.Fatalf("unexpected error creating display: %v", err)
	}

	result := display.Twilight()

	// Should handle gracefully, not return empty or panic
	if result == "" {
		t.Error("Twilight() returned empty string, expected graceful handling")
	}
	if !strings.Contains(result, "No") && !strings.Contains(result, "data") {
		// If it doesn't show a "no data" message, that's also acceptable
		// as long as it doesn't panic or return empty
		t.Logf("Twilight() = %q", result)
	}
}

func TestWarnings_MultipleAlerts(t *testing.T) {
	tests := []struct {
		name   string
		alerts []api.Alert
		want   string
	}{
		{
			name:   "no alerts",
			alerts: []api.Alert{},
			want:   "None",
		},
		{
			name:   "one alert",
			alerts: []api.Alert{{Event: "Flood Warning"}},
			want:   "Flood Warning",
		},
		{
			name: "multiple alerts",
			alerts: []api.Alert{
				{Event: "Flood Warning"},
				{Event: "Wind Advisory"},
			},
			want: "Flood Warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &api.Response{
				Forecast: api.Forecast{Forecastday: []api.ForecastDay{{}}},
				Alerts:   api.Alerts{Alert: tt.alerts},
			}

			display, err := NewDisplay(data, true)
			if err != nil {
				t.Fatalf("unexpected error creating display: %v", err)
			}

			result := display.Warnings()

			if !strings.Contains(result, tt.want) {
				t.Errorf("Warnings() = %q, want string containing %q", result, tt.want)
			}

			// For multiple alerts, verify second alert is also present
			if tt.name == "multiple alerts" {
				if !strings.Contains(result, "Wind Advisory") {
					t.Errorf("Warnings() = %q, want string containing 'Wind Advisory'", result)
				}
			}
		})
	}
}

func TestDailyForecast_NoData(t *testing.T) {
	// Only one day (today) means no future days to show
	data := &api.Response{
		Forecast: api.Forecast{
			Forecastday: []api.ForecastDay{{}},
		},
	}

	display, err := NewDisplay(data, true)
	if err != nil {
		t.Fatalf("unexpected error creating display: %v", err)
	}

	result := display.DailyForecast()

	// Should indicate no data available when only today exists
	if !strings.Contains(result, "No data available") {
		t.Errorf("DailyForecast() = %q, want string containing 'No data available'", result)
	}
}

func TestDailyForecast_MultipleDays(t *testing.T) {
	data := &api.Response{
		Forecast: api.Forecast{
			Forecastday: []api.ForecastDay{
				{
					// Today - should be skipped
					Date: "2025-12-01",
					Day: api.Day{
						MaxTempC:     15,
						MinTempC:     8,
						ChanceOfRain: 20,
						Condition:    api.Condition{Text: "Sunny"},
					},
				},
				{
					// Tomorrow - should be shown
					Date: "2025-12-02",
					Day: api.Day{
						MaxTempC:     12,
						MinTempC:     5,
						ChanceOfRain: 60,
						Condition:    api.Condition{Text: "Light rain"},
					},
				},
				{
					// Day after - should be shown
					Date: "2025-12-03",
					Day: api.Day{
						MaxTempC:     10,
						MinTempC:     3,
						ChanceOfRain: 30,
						Condition:    api.Condition{Text: "Cloudy"},
					},
				},
			},
		},
	}

	display, err := NewDisplay(data, true)
	if err != nil {
		t.Fatalf("unexpected error creating display: %v", err)
	}

	result := display.DailyForecast()

	// Should NOT contain today (Mon 01)
	if strings.Contains(result, "Mon 01") {
		t.Errorf("DailyForecast() = %q, should NOT contain today 'Mon 01'", result)
	}

	// Should contain future days
	if !strings.Contains(result, "Tue 02") {
		t.Errorf("DailyForecast() = %q, want string containing 'Tue 02'", result)
	}
	if !strings.Contains(result, "Wed 03") {
		t.Errorf("DailyForecast() = %q, want string containing 'Wed 03'", result)
	}

	// Should contain temperatures for tomorrow (now in separate columns)
	if !strings.Contains(result, "12°C") {
		t.Errorf("DailyForecast() = %q, want string containing '12°C'", result)
	}
	if !strings.Contains(result, "5°C") {
		t.Errorf("DailyForecast() = %q, want string containing '5°C'", result)
	}

	// Should contain column headers
	if !strings.Contains(result, "High") {
		t.Errorf("DailyForecast() = %q, want string containing 'High'", result)
	}
	if !strings.Contains(result, "Low") {
		t.Errorf("DailyForecast() = %q, want string containing 'Low'", result)
	}

	// Should contain rain chance
	if !strings.Contains(result, "60%") {
		t.Errorf("DailyForecast() = %q, want string containing '60%%'", result)
	}

	// Should contain condition text
	if !strings.Contains(result, "Light rain") {
		t.Errorf("DailyForecast() = %q, want string containing 'Light rain'", result)
	}
}
