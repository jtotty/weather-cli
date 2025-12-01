package config

import (
	"testing"
)

func TestNew_RequiresAPIKey(t *testing.T) {
	tests := []struct {
		name    string
		envVal  string
		wantErr bool
	}{
		{"missing key", "", true},
		{"valid key", "abc123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("WEATHER_API_KEY", tt.envVal)

			cfg, err := New()

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if cfg != nil {
					t.Error("expected nil config when error occurs")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if cfg == nil {
					t.Error("expected config, got nil")
				}
				if cfg != nil && cfg.APIKey != tt.envVal {
					t.Errorf("APIKey = %q, want %q", cfg.APIKey, tt.envVal)
				}
			}
		})
	}
}

func TestNew_DefaultValues(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "test-key")

	cfg, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Location != "auto:ip" {
		t.Errorf("Location = %q, want %q", cfg.Location, "auto:ip")
	}
	if cfg.Days != 1 {
		t.Errorf("Days = %d, want %d", cfg.Days, 1)
	}
	if cfg.IncludeAQI != true {
		t.Errorf("IncludeAQI = %v, want %v", cfg.IncludeAQI, true)
	}
	if cfg.Alerts != true {
		t.Errorf("Alerts = %v, want %v", cfg.Alerts, true)
	}
	if cfg.IsLocal != true {
		t.Errorf("IsLocal = %v, want %v", cfg.IsLocal, true)
	}
}

func TestSetLocation_SetsIsLocalFalse(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "test-key")

	cfg, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.IsLocal {
		t.Error("IsLocal should be true before SetLocation")
	}

	cfg.SetLocation("Paris")

	if cfg.Location != "Paris" {
		t.Errorf("Location = %q, want %q", cfg.Location, "Paris")
	}
	if cfg.IsLocal {
		t.Error("IsLocal should be false after SetLocation")
	}
}
