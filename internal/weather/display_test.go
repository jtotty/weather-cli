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
