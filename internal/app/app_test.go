package app

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/config"
)

type mockWeatherService struct {
	response *weather.Response
	err      error
}

func (m *mockWeatherService) GetWeather(_ context.Context) (*weather.Response, error) {
	return m.response, m.err
}

func validWeatherResponse() *weather.Response {
	return &weather.Response{
		Location: weather.Location{
			Name:      "London",
			Country:   "United Kingdom",
			LocalTime: "2024-01-15 14:30",
		},
		Current: weather.Current{
			TempC:         12.5,
			FeelsLike:     10.2,
			Humidity:      75,
			WindSpeed:     8.5,
			WindDirection: "SW",
			Condition:     weather.Condition{Text: "Partly cloudy"},
		},
		Forecast: weather.Forecast{
			Forecastday: []weather.ForecastDay{
				{
					Date: "2024-01-15",
					Day: weather.Day{
						MaxTempC:     14.0,
						MinTempC:     8.0,
						Condition:    weather.Condition{Text: "Cloudy"},
						ChanceOfRain: 30,
					},
					Hour: []weather.Hour{
						{TimeEpoch: 1705320000, TempC: 10.0, Condition: weather.Condition{Text: "Clear"}},
					},
					Astro: weather.Astro{Sunrise: "07:45 AM", Sunset: "04:30 PM"},
				},
			},
		},
	}
}

func TestApp_Run(t *testing.T) {
	tests := []struct {
		name            string
		location        string
		mockResponse    *weather.Response
		mockErr         error
		wantErr         bool
		wantInOutput    []string
		wantNotInOutput []string
	}{
		{
			name:         "successful weather fetch",
			location:     "London",
			mockResponse: validWeatherResponse(),
			wantErr:      false,
			wantInOutput: []string{"London", "United Kingdom"},
		},
		{
			name:         "empty location uses config default",
			location:     "",
			mockResponse: validWeatherResponse(),
			wantErr:      false,
			wantInOutput: []string{"London"},
		},
		{
			name:     "service error propagates",
			location: "InvalidCity",
			mockErr:  errors.New("API error: city not found"),
			wantErr:  true,
		},
		{
			name:     "context canceled error",
			location: "London",
			mockErr:  context.Canceled,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.NewWithAPIKey("test-api-key")
			mockSvc := &mockWeatherService{
				response: tt.mockResponse,
				err:      tt.mockErr,
			}
			var buf bytes.Buffer

			application := NewApp(cfg, mockSvc, &buf)
			err := application.Run(context.Background(), tt.location)

			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			output := buf.String()
			for _, want := range tt.wantInOutput {
				if !strings.Contains(output, want) {
					t.Errorf("Run() output should contain %q, got %q", want, output)
				}
			}
		})
	}
}

func TestApp_Run_LocationSetsIsLocalFalse(t *testing.T) {
	cfg := config.NewWithAPIKey("test-api-key")
	if !cfg.IsLocal {
		t.Fatal("expected IsLocal to be true initially")
	}

	mockSvc := &mockWeatherService{response: validWeatherResponse()}
	var buf bytes.Buffer

	application := NewApp(cfg, mockSvc, &buf)
	err := application.Run(context.Background(), "Paris")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.IsLocal {
		t.Error("expected IsLocal to be false after setting location")
	}
	if cfg.Location != "Paris" {
		t.Errorf("expected location to be Paris, got %s", cfg.Location)
	}
}

func TestApp_Run_ContextCancellation(t *testing.T) {
	cfg := config.NewWithAPIKey("test-api-key")
	mockSvc := &mockWeatherService{err: context.Canceled}
	var buf bytes.Buffer

	application := NewApp(cfg, mockSvc, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := application.Run(ctx, "London")

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestNewApp(t *testing.T) {
	cfg := config.NewWithAPIKey("test-api-key")
	mockSvc := &mockWeatherService{response: validWeatherResponse()}
	var buf bytes.Buffer

	application := NewApp(cfg, mockSvc, &buf)

	if application.config != cfg {
		t.Error("expected config to be set")
	}
	if application.service != mockSvc {
		t.Error("expected service to be set")
	}
	if application.output != &buf {
		t.Error("expected output to be set")
	}
}
