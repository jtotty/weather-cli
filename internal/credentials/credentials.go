// Package credentials owns where the API key lives: retrieval, storage,
// deletion, availability, and precedence between sources.
package credentials

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	serviceName = "weather-cli"
	apiKeyName  = "api-key"
	envVarName  = "WEATHER_API_KEY"
)

var (
	ErrNoAPIKey           = errors.New("no API key configured")
	ErrKeyringUnavailable = errors.New("OS keyring unavailable")
)

// Store is the credential store seam: the single answer to where the API key
// comes from.
type Store interface {
	Get() (string, error)
	Set(key string) error
	Delete() error
	Available() bool
}

// Keyring is the production Store adapter backed by the OS keyring.
type Keyring struct{}

func NewKeyring() *Keyring {
	return &Keyring{}
}

func (k *Keyring) Get() (string, error) {
	key, err := keyring.Get(serviceName, apiKeyName)
	if err == nil && key != "" {
		return key, nil
	}

	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNoAPIKey
	}

	if errors.Is(err, keyring.ErrUnsupportedPlatform) {
		return "", ErrKeyringUnavailable
	}

	if err != nil {
		return "", fmt.Errorf("failed to retrieve API key: %w", err)
	}

	return "", ErrNoAPIKey
}

// normalizeKey trims a key and rejects blank ones, so every Store adapter
// shares the same write semantics.
func normalizeKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", errors.New("API key cannot be empty")
	}
	return key, nil
}

func (k *Keyring) Set(key string) error {
	key, err := normalizeKey(key)
	if err != nil {
		return err
	}

	if err := keyring.Set(serviceName, apiKeyName, key); err != nil {
		if errors.Is(err, keyring.ErrUnsupportedPlatform) {
			return ErrKeyringUnavailable
		}
		return fmt.Errorf("failed to store API key: %w", err)
	}

	return nil
}

func (k *Keyring) Delete() error {
	if err := keyring.Delete(serviceName, apiKeyName); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

func (k *Keyring) Available() bool {
	_, err := keyring.Get(serviceName, "test-availability")
	if err == nil {
		return true
	}

	return errors.Is(err, keyring.ErrNotFound)
}

// EnvOverride decorates a Store so WEATHER_API_KEY beats the stored key on
// reads. Set, Delete, and Available pass through to the inner store.
type EnvOverride struct {
	inner Store
}

func NewEnvOverride(inner Store) *EnvOverride {
	return &EnvOverride{inner: inner}
}

func (e *EnvOverride) Get() (string, error) {
	if key := os.Getenv(envVarName); key != "" {
		return key, nil
	}
	return e.inner.Get()
}

func (e *EnvOverride) Set(key string) error {
	return e.inner.Set(key)
}

func (e *EnvOverride) Delete() error {
	return e.inner.Delete()
}

func (e *EnvOverride) Available() bool {
	return e.inner.Available()
}
