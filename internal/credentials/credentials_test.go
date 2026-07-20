package credentials

import (
	"errors"
	"testing"
)

const testKey = "my-key"

func TestFake_GetReturnsErrNoAPIKeyWhenEmpty(t *testing.T) {
	fake := &Fake{}

	_, err := fake.Get()
	if !errors.Is(err, ErrNoAPIKey) {
		t.Errorf("Get() error = %v, want ErrNoAPIKey", err)
	}
}

func TestFake_SetThenGetRoundTrips(t *testing.T) {
	fake := &Fake{}

	if err := fake.Set(testKey); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	key, err := fake.Get()
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if key != testKey {
		t.Errorf("Get() = %q, want %q", key, testKey)
	}
}

func TestFake_SetTrimsWhitespace(t *testing.T) {
	fake := &Fake{}

	if err := fake.Set("  my-key  "); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	key, _ := fake.Get()
	if key != testKey {
		t.Errorf("Get() = %q, want %q", key, testKey)
	}
}

func TestFake_SetRejectsEmptyKey(t *testing.T) {
	fake := &Fake{}

	if err := fake.Set("   "); err == nil {
		t.Error("Set() with blank key should return an error")
	}
}

func TestFake_DeleteRemovesKey(t *testing.T) {
	fake := &Fake{}
	if err := fake.Set(testKey); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if err := fake.Delete(); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := fake.Get()
	if !errors.Is(err, ErrNoAPIKey) {
		t.Errorf("Get() after Delete() error = %v, want ErrNoAPIKey", err)
	}
}

func TestFake_DeleteWithoutKeyIsNoError(t *testing.T) {
	fake := &Fake{}

	if err := fake.Delete(); err != nil {
		t.Errorf("Delete() on empty store error = %v, want nil", err)
	}
}

func TestFake_UnavailableFailsEveryOperation(t *testing.T) {
	fake := &Fake{Key: testKey, Unavailable: true}

	if _, err := fake.Get(); !errors.Is(err, ErrKeyringUnavailable) {
		t.Errorf("Get() error = %v, want ErrKeyringUnavailable", err)
	}
	if err := fake.Set("other-key"); !errors.Is(err, ErrKeyringUnavailable) {
		t.Errorf("Set() error = %v, want ErrKeyringUnavailable", err)
	}
	if err := fake.Delete(); !errors.Is(err, ErrKeyringUnavailable) {
		t.Errorf("Delete() error = %v, want ErrKeyringUnavailable", err)
	}
}

func TestFake_Available(t *testing.T) {
	fake := &Fake{}
	if !fake.Available() {
		t.Error("Available() = false, want true by default")
	}

	fake.Unavailable = true
	if fake.Available() {
		t.Error("Available() = true, want false when Unavailable is set")
	}
}

func TestEnvOverride_EnvVarBeatsStoredKey(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "env-key")

	fake := &Fake{}
	if err := fake.Set("stored-key"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	store := NewEnvOverride(fake)

	key, err := store.Get()
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if key != "env-key" {
		t.Errorf("Get() = %q, want %q", key, "env-key")
	}
}

func TestEnvOverride_FallsBackToInnerStore(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "")

	fake := &Fake{}
	if err := fake.Set("stored-key"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	store := NewEnvOverride(fake)

	key, err := store.Get()
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if key != "stored-key" {
		t.Errorf("Get() = %q, want %q", key, "stored-key")
	}
}

func TestEnvOverride_SetDeleteAvailablePassThrough(t *testing.T) {
	t.Setenv("WEATHER_API_KEY", "env-key")

	fake := &Fake{}
	store := NewEnvOverride(fake)

	if err := store.Set("stored-key"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if key, err := fake.Get(); err != nil || key != "stored-key" {
		t.Errorf("inner store Get() = %q, %v, want %q", key, err, "stored-key")
	}

	if !store.Available() {
		t.Error("Available() = false, want passthrough true")
	}

	if err := store.Delete(); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := fake.Get(); !errors.Is(err, ErrNoAPIKey) {
		t.Errorf("inner store Get() after Delete() error = %v, want ErrNoAPIKey", err)
	}
}
