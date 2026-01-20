# Weather CLI

A command-line tool for checking weather forecasts, built with Go.

## Features

- Current conditions with temperature, humidity, wind, and air quality index
- 7-day forecast with daily high/low temperatures
- 24-hour hourly forecast
- Sunrise and sunset times
- Weather alerts and warnings
- Automatic location detection via IP geolocation
- Supports city names, ZIP codes, and coordinates
- 30-minute response caching for faster repeat queries
- Secure API key storage using OS keyring

## Requirements

- Go 1.25 or later
- A free API key from [WeatherAPI.com](https://www.weatherapi.com/)

## Installation

```bash
go install github.com/jtotty/weather-cli@latest
```

Or build from source:

```bash
git clone https://github.com/jtotty/weather-cli.git
cd weather-cli
go build -o weather-cli .
```

## Configuration

Get a free API key from [weatherapi.com](https://www.weatherapi.com/) and configure it using one of these methods:

### Option 1: Interactive Setup (Recommended)

```bash
weather-cli --setup
```

This stores your API key securely in your OS keyring (macOS Keychain, Linux Secret Service, or Windows Credential Manager).

### Option 2: Environment Variable

```bash
export WEATHER_API_KEY=your_api_key_here
```

## Usage

```bash
# Weather for your current location (via IP geolocation)
weather-cli

# Weather for a specific city
weather-cli London

# City names with spaces (use quotes)
weather-cli "New York"

# ZIP code
weather-cli 10001

# Coordinates (latitude,longitude)
weather-cli 51.5,-0.1
```

### Options

| Flag | Description |
|------|-------------|
| `-h`, `--help` | Show help message |
| `-v`, `--version` | Show version |
| `--setup` | Configure API key |
| `--delete-key` | Remove stored API key |

## Development

```bash
# Run
go run . [location]

# Test
go test ./...

# Test with race detection
go test -race ./...

# Lint
golangci-lint run

# Build
go build -o weather-cli .
```

## License

MIT License - see [LICENSE](LICENSE) for details.
