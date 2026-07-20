package cli

import (
	"errors"
	"testing"

	"github.com/jtotty/weather-cli/internal/credentials"
)

const testKey = "my-key"

func fakeInput(key string) KeyReader {
	return func() (string, error) { return key, nil }
}

func TestRunSetup_StoresKeyFromInput(t *testing.T) {
	store := &credentials.Fake{}

	if err := RunSetup(store, fakeInput(testKey)); err != nil {
		t.Fatalf("RunSetup() error = %v", err)
	}

	key, err := store.Get()
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if key != testKey {
		t.Errorf("stored key = %q, want %q", key, testKey)
	}
}

func TestRunSetup_ErrorsWhenStoreUnavailable(t *testing.T) {
	store := &credentials.Fake{Unavailable: true}

	err := RunSetup(store, fakeInput(testKey))
	if err == nil {
		t.Fatal("RunSetup() should error when the store is unavailable")
	}

	store.Unavailable = false
	if _, getErr := store.Get(); !errors.Is(getErr, credentials.ErrNoAPIKey) {
		t.Error("no key should be stored when the store is unavailable")
	}
}

func TestRunSetup_PropagatesInputError(t *testing.T) {
	store := &credentials.Fake{}
	inputErr := errors.New("read failed")
	input := func() (string, error) { return "", inputErr }

	if err := RunSetup(store, input); !errors.Is(err, inputErr) {
		t.Errorf("RunSetup() error = %v, want %v", err, inputErr)
	}
}

func TestRunDeleteKey_RemovesStoredKey(t *testing.T) {
	store := &credentials.Fake{Key: testKey}

	if err := RunDeleteKey(store); err != nil {
		t.Fatalf("RunDeleteKey() error = %v", err)
	}

	if _, err := store.Get(); !errors.Is(err, credentials.ErrNoAPIKey) {
		t.Errorf("Get() after RunDeleteKey() error = %v, want ErrNoAPIKey", err)
	}
}

func TestLoadConfig_ReturnsConfigWhenKeyPresent(t *testing.T) {
	store := &credentials.Fake{Key: testKey}

	cfg, err := LoadConfig(store, fakeInput("unused"))
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.APIKey != testKey {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, testKey)
	}
}

func TestLoadConfig_MissingKeyRunsSetup(t *testing.T) {
	store := &credentials.Fake{}

	cfg, err := LoadConfig(store, fakeInput("new-key"))
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg.APIKey != "new-key" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "new-key")
	}
}

func TestLoadConfig_MissingKeyAndFailedSetupReturnsError(t *testing.T) {
	store := &credentials.Fake{}
	inputErr := errors.New("read failed")
	input := func() (string, error) { return "", inputErr }

	if _, err := LoadConfig(store, input); !errors.Is(err, inputErr) {
		t.Errorf("LoadConfig() error = %v, want %v", err, inputErr)
	}
}
