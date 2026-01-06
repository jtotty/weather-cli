// Package credentials provides secure storage for API keys using the OS keyring.
package credentials

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/zalando/go-keyring"
	"golang.org/x/term"
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

func GetAPIKey() (string, error) {
	if key := os.Getenv(envVarName); key != "" {
		return key, nil
	}

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

func SetAPIKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("API key cannot be empty")
	}

	if err := keyring.Set(serviceName, apiKeyName, key); err != nil {
		if errors.Is(err, keyring.ErrUnsupportedPlatform) {
			return ErrKeyringUnavailable
		}
		return fmt.Errorf("failed to store API key: %w", err)
	}

	return nil
}

func DeleteAPIKey() error {
	if err := keyring.Delete(serviceName, apiKeyName); err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to delete API key: %w", err)
	}
	return nil
}

func PromptForAPIKey() (string, error) {
	fmt.Print("Enter your Weather API key: ")

	keyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()

	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	key := strings.TrimSpace(string(keyBytes))
	if key == "" {
		return "", errors.New("API key cannot be empty")
	}

	return key, nil
}

func IsKeyringAvailable() bool {
	_, err := keyring.Get(serviceName, "test-availability")
	if err == nil {
		return true
	}

	if errors.Is(err, keyring.ErrNotFound) {
		return true
	}

	return false
}

func HasStoredAPIKey() bool {
	key, err := keyring.Get(serviceName, apiKeyName)
	return err == nil && key != ""
}
