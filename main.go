package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/cache"
	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/credentials"
	"github.com/jtotty/weather-cli/internal/ui"
	weatherdisplay "github.com/jtotty/weather-cli/internal/weather"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	if handleFlags() {
		return
	}

	cfg, err := config.New()
	if err != nil {
		// If no API key configured, prompt for setup
		if errors.Is(err, credentials.ErrNoAPIKey) {
			fmt.Println("No API key configured.")
			fmt.Println()
			if err := runSetup(); err != nil {
				exitWithErrorf("Setup failed: %v", err)
			}
			// Retry loading config after setup
			cfg, err = config.New()
			if err != nil {
				exitWithErrorf("Error loading config: %v", err)
			}
		} else {
			exitWithErrorf("Error loading config: %v", err)
		}
	}

	if err := parseLocationArg(cfg); err != nil {
		exitWithErrorf("%v", err)
	}

	data, err := fetchWeatherData(cfg)
	if err != nil {
		exitWithErrorf("Error fetching weather: %v", err)
	}

	if err := displayWeather(data, cfg.IsLocal); err != nil {
		exitWithErrorf("Error creating display: %v", err)
	}
}

// handleFlags checks for command line flags and handles them.
// Returns true if a flag was handled and the program should exit.
func handleFlags() bool {
	if len(os.Args) < 2 {
		return false
	}

	arg := os.Args[1]

	switch arg {
	case "--version", "-v":
		fmt.Printf("weather-cli %s\n", version)
		return true

	case "--help", "-h":
		printHelp()
		return true

	case "--setup":
		if err := runSetup(); err != nil {
			exitWithErrorf("Setup failed: %v", err)
		}
		return true

	case "--delete-key":
		if err := runDeleteKey(); err != nil {
			exitWithErrorf("Failed to delete API key: %v", err)
		}
		fmt.Println("API key deleted from keyring.")
		return true
	}

	return false
}

// printHelp prints usage information.
func printHelp() {
	fmt.Printf(`weather-cli %s

USAGE:
    weather-cli [OPTIONS] [LOCATION]

ARGUMENTS:
    [LOCATION]    Location for weather lookup (city name, zip code, coordinates)
                  If omitted, uses your current location via IP geolocation

OPTIONS:
    -h, --help        Show this help message
    -v, --version     Show version information
    --setup           Configure your Weather API key (stored in OS keyring)
    --delete-key      Remove stored API key from OS keyring

EXAMPLES:
    weather-cli                     # Weather for current location
    weather-cli London              # Weather for London
    weather-cli "New York"          # Weather for New York (use quotes for spaces)
    weather-cli 10001               # Weather for ZIP code 10001
    weather-cli 51.5,-0.1           # Weather for coordinates

API KEY:
    Get a free API key from https://www.weatherapi.com/
    Run 'weather-cli --setup' to configure your API key.
    
    Alternatively, set the WEATHER_API_KEY environment variable.
`, version)
}

// runSetup prompts the user for their API key and stores it securely.
func runSetup() error {
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

// runDeleteKey removes the stored API key from the keyring.
func runDeleteKey() error {
	return credentials.DeleteAPIKey()
}

// parseLocationArg parses the location from command line arguments.
func parseLocationArg(cfg *config.Config) error {
	if len(os.Args) < 2 {
		return nil
	}

	location := strings.TrimSpace(os.Args[1])

	// Skip if it's a flag
	if strings.HasPrefix(location, "-") {
		return nil
	}

	if location == "" {
		return fmt.Errorf("location cannot be empty")
	}

	if len(location) > 100 {
		return fmt.Errorf("location too long (max 100 characters)")
	}

	cfg.SetLocation(location)
	return nil
}

// fetchWeatherData retrieves weather data from cache or API.
func fetchWeatherData(cfg *config.Config) (*weather.Response, error) {
	weatherCache, err := cache.New(cache.DefaultTTL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: cache unavailable: %v\n", err)
	}

	// Try cache first
	if weatherCache != nil {
		if data := weatherCache.Get(cfg.Location); data != nil {
			return data, nil
		}
	}

	// Fetch from API
	data, err := fetchFromAPI(cfg)
	if err != nil {
		return nil, err
	}

	// Cache the response
	if weatherCache != nil {
		if cacheErr := weatherCache.Set(cfg.Location, data); cacheErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to cache data: %v\n", cacheErr)
		}
	}

	return data, nil
}

// fetchFromAPI fetches weather data from the weather API.
func fetchFromAPI(cfg *config.Config) (*weather.Response, error) {
	client := weather.NewClient(cfg.APIKey)
	ctx := context.Background()

	return client.Fetch(ctx, weather.FetchOptions{
		Location:   cfg.Location,
		Days:       cfg.Days,
		IncludeAQI: cfg.IncludeAQI,
		Alerts:     cfg.Alerts,
	})
}

// displayWeather renders the weather data to stdout.
func displayWeather(data *weather.Response, isLocal bool) error {
	display, err := weatherdisplay.NewDisplay(data, isLocal)
	if err != nil {
		return err
	}

	fmt.Print(display.Heading())
	ui.Spacer()

	fmt.Print(display.Time())
	ui.Spacer()

	fmt.Print(display.CurrentConditions())
	ui.Spacer()

	fmt.Print(display.HourlyForecast())
	ui.Spacer()

	fmt.Print(display.DailyForecast())
	ui.Spacer()

	fmt.Print(display.Twilight())
	ui.Spacer()

	fmt.Print(display.Warnings())

	return nil
}

// exitWithErrorf prints an error message to stderr and exits with code 1.
func exitWithErrorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
