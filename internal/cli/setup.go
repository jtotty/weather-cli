package cli

import (
	"fmt"

	"github.com/jtotty/weather-cli/internal/credentials"
)

func RunSetup() error {
	if !credentials.IsKeyringAvailable() {
		return fmt.Errorf("OS keyring is not available on this system. Set WEATHER_API_KEY environment variable instead")
	}

	fmt.Println("Get a free API key from https://www.weatherapi.com/")
	fmt.Println()

	key, err := credentials.PromptForAPIKey()
	if err != nil {
		return err
	}

	if err := credentials.SetAPIKey(key); err != nil {
		return err
	}

	fmt.Println("API key stored securely in OS keyring.")
	return nil
}

func RunDeleteKey() error {
	if err := credentials.DeleteAPIKey(); err != nil {
		return err
	}
	fmt.Println("API key deleted from keyring.")
	return nil
}
