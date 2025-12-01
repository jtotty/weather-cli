package weather

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildURL_EncodesSpecialCharacters(t *testing.T) {
	client := NewClient("test-key")

	tests := []struct {
		name     string
		location string
	}{
		{"ampersand", "London & Paris"},
		{"equals", "key=value"},
		{"space", "New York"},
		{"question mark", "where?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := FetchOptions{Location: tt.location, Days: 1}
			url := client.BuildURL(opts)

			// The location should be URL-encoded, not appear as raw special chars
			// after the q= parameter
			if strings.Contains(url, "q="+tt.location) {
				t.Errorf("URL contains unencoded location: %s", url)
			}

			// URL should still be valid and contain the encoded location
			if !strings.Contains(url, "q=") {
				t.Errorf("URL missing q parameter: %s", url)
			}
		})
	}
}

func TestBuildURL_IncludesOptionalParams(t *testing.T) {
	client := NewClient("test-key")

	tests := []struct {
		name       string
		opts       FetchOptions
		wantParams []string
		dontWant   []string
	}{
		{
			name: "all options enabled",
			opts: FetchOptions{
				Location:   "London",
				Days:       3,
				IncludeAQI: true,
				Alerts:     true,
			},
			wantParams: []string{"aqi=yes", "alerts=yes", "days=3"},
		},
		{
			name: "options disabled",
			opts: FetchOptions{
				Location:   "London",
				Days:       1,
				IncludeAQI: false,
				Alerts:     false,
			},
			dontWant: []string{"aqi=", "alerts="},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := client.BuildURL(tt.opts)

			for _, param := range tt.wantParams {
				if !strings.Contains(url, param) {
					t.Errorf("URL = %q, want param %q", url, param)
				}
			}

			for _, param := range tt.dontWant {
				if strings.Contains(url, param) {
					t.Errorf("URL = %q, should not contain %q", url, param)
				}
			}
		})
	}
}

func TestFetch_Success(t *testing.T) {
	mockResponse := Response{
		Location: Location{
			Name:      "London",
			Country:   "UK",
			LocalTime: "2024-01-15 12:00",
		},
		Current: Current{
			TempC: 15,
		},
		Forecast: Forecast{
			Forecastday: []ForecastDay{
				{
					Astro: Astro{Sunrise: "07:00 AM", Sunset: "05:00 PM"},
					Hour:  []Hour{},
				},
			},
		},
		Alerts: Alerts{Alert: []Alert{}},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request has expected parameters
		if !strings.Contains(r.URL.RawQuery, "key=test-key") {
			t.Error("request missing API key")
		}
		if !strings.Contains(r.URL.RawQuery, "q=London") {
			t.Error("request missing location")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := NewTestClient("test-key", server.URL)
	ctx := context.Background()

	response, err := client.Fetch(ctx, FetchOptions{
		Location: "London",
		Days:     1,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Location.Name != "London" {
		t.Errorf("Location.Name = %q, want %q", response.Location.Name, "London")
	}

	if response.Location.Country != "UK" {
		t.Errorf("Location.Country = %q, want %q", response.Location.Country, "UK")
	}

	if response.Current.TempC != 15 {
		t.Errorf("Current.TempC = %v, want %v", response.Current.TempC, 15)
	}
}

func TestFetch_HTTPErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    string
	}{
		{"unauthorized", http.StatusUnauthorized, "401"},
		{"not found", http.StatusNotFound, "404"},
		{"server error", http.StatusInternalServerError, "500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			client := NewTestClient("test-key", server.URL)
			ctx := context.Background()

			_, err := client.Fetch(ctx, FetchOptions{Location: "London", Days: 1})

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want error containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}
