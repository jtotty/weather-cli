package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/credentials"
)

// KeyReader reads an API key from the user.
type KeyReader func() (string, error)

// ReadKeyFromTerminal is the production KeyReader: it prompts on stdout and
// reads the key from the terminal without echoing it.
func ReadKeyFromTerminal() (string, error) {
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

// RunSetup is the setup wizard: it checks store availability, prompts for a
// key via input, and stores it.
func RunSetup(store credentials.Store, input KeyReader) error {
	if !store.Available() {
		return fmt.Errorf("OS keyring is not available on this system. Set WEATHER_API_KEY environment variable instead")
	}

	fmt.Println("Get a free API key from https://www.weatherapi.com/")
	fmt.Println()

	key, err := input()
	if err != nil {
		return err
	}

	if err := store.Set(key); err != nil {
		return err
	}

	fmt.Println("API key stored securely in OS keyring.")
	return nil
}

func RunDeleteKey(store credentials.Store) error {
	if err := store.Delete(); err != nil {
		return err
	}
	fmt.Println("API key deleted from keyring.")
	return nil
}

// LoadConfig loads the app config, guiding the user through setup when no
// API key is configured yet.
func LoadConfig(store credentials.Store, input KeyReader) (*config.Config, error) {
	cfg, err := config.New(store)
	if err == nil {
		return cfg, nil
	}

	if errors.Is(err, credentials.ErrNoAPIKey) {
		fmt.Println("No API key configured.")
		fmt.Println()
		if setupErr := RunSetup(store, input); setupErr != nil {
			return nil, fmt.Errorf("setup failed: %w", setupErr)
		}
		return config.New(store)
	}

	return nil, fmt.Errorf("error loading config: %w", err)
}
