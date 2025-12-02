# AGENTS.md

## Build & Run Commands
- Build: `go build -o weather-cli .`
- Run: `go run .` or `go run . <location>`
- Test all: `go test ./...`
- Test single: `go test -run TestName ./path/to/package`
- Lint: `go vet ./...`

## API Key Setup
- Run `weather-cli --setup` to configure your API key (stored in OS keyring)
- Alternatively, set the `WEATHER_API_KEY` environment variable
- Get a free API key from https://www.weatherapi.com/

## Code Style Guidelines
- **Imports**: Group stdlib, then external, then internal (`github.com/jtotty/weather-cli/internal/...`)
- **Formatting**: Use `gofmt` or `goimports`; tabs for indentation
- **Naming**: PascalCase for exports, camelCase for private; descriptive names
- **Error handling**: Return errors with context using `fmt.Errorf("action failed: %w", err)`
- **Types**: Use structs with JSON tags for API responses; prefer `float32` for numeric weather data
- **Strings**: Use `strings.Builder` for concatenation in methods

## Project Structure
- `main.go` - Entry point, handles CLI flags, initializes weather service
- `internal/config/` - Configuration and API URL building
- `internal/credentials/` - Secure API key storage (OS keyring)
- `internal/cache/` - File-based weather data caching (30min TTL)
- `internal/api/weather/` - Weather API client
- `internal/weather/` - Weather display formatting
- `internal/ui/` - UI helpers, icons, emoji mappings
