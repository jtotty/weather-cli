// Package cli handles command-line argument parsing and output.
package cli

import (
	"fmt"
	"os"
	"strings"
)

type CommandType int

const (
	CommandWeather CommandType = iota
	CommandHelp
	CommandVersion
	CommandSetup
	CommandDeleteKey
)

type Command struct {
	Type     CommandType
	Location string
}

func Parse(args []string) Command {
	if len(args) < 2 {
		return Command{Type: CommandWeather}
	}

	arg := args[1]

	switch arg {
	case "--help", "-h":
		return Command{Type: CommandHelp}
	case "--version", "-v":
		return Command{Type: CommandVersion}
	case "--setup":
		return Command{Type: CommandSetup}
	case "--delete-key":
		return Command{Type: CommandDeleteKey}
	default:
		// Treat as location if not a flag
		if strings.HasPrefix(arg, "-") {
			return Command{Type: CommandHelp}
		}
		return Command{Type: CommandWeather, Location: arg}
	}
}

func PrintHelp(version string) {
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

func PrintVersion(version string) {
	fmt.Printf("weather-cli %s\n", version)
}

func ExitWithError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
