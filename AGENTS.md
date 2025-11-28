# AGENTS.md

## Build & Run Commands
- Build: `go build -o weather-cli .`
- Run: `go run .` or `go run . <location>`
- Test all: `go test ./...`
- Test single: `go test -run TestName ./path/to/package`
- Lint: `go vet ./...`

## Code Style Guidelines
- **Imports**: Group stdlib, then external, then internal (`github.com/jtotty/weather-cli/internal/...`)
- **Formatting**: Use `gofmt` or `goimports`; tabs for indentation
- **Naming**: PascalCase for exports, camelCase for private; descriptive names
- **Error handling**: Return errors with context using `fmt.Errorf("action failed: %w", err)`
- **Types**: Use structs with JSON tags for API responses; prefer `float32` for numeric weather data
- **Strings**: Use `strings.Builder` for concatenation in methods

## Project Structure
- `main.go` - Entry point, loads .env, initializes weather service
- `internal/config/` - Configuration and API URL building
- `internal/weather/` - Weather data fetching and formatting
- `internal/ui/` - UI helpers, icons, emoji mappings
