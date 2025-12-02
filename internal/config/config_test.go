package config

import (
	"testing"
)

func TestNew_WithEnvAPIKey(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "test-api-key")

	cfg, err := New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected config, got nil")
	}
	if cfg.APIKey != "test-api-key" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "test-api-key")
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
	if cfg.Days != 7 {
		t.Errorf("Days = %d, want %d", cfg.Days, 7)
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
