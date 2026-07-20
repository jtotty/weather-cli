package config

import (
	"errors"
	"testing"

	"github.com/jtotty/weather-cli/internal/credentials"
)

func TestNew_UsesKeyFromStore(t *testing.T) {
	store := &credentials.Fake{Key: "test-api-key"}

	cfg, err := New(store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.APIKey != "test-api-key" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "test-api-key")
	}
}

func TestNew_MissingKeyReturnsErrNoAPIKey(t *testing.T) {
	store := &credentials.Fake{}

	_, err := New(store)
	if !errors.Is(err, credentials.ErrNoAPIKey) {
		t.Errorf("New() error = %v, want ErrNoAPIKey", err)
	}
}

func TestNew_DefaultValues(t *testing.T) {
	store := &credentials.Fake{Key: "test-key"}

	cfg, err := New(store)
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
	store := &credentials.Fake{Key: "test-key"}

	cfg, err := New(store)
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
