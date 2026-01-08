# AGENTS.md

## Overview
CLI weather tool using WeatherAPI.com. Go 1.25, single binary, OS keyring for credentials.

## Structure
```
.
├── main.go                      # Entry point, CLI dispatch
├── internal/
│   ├── api/weather/             # HTTP client, API types
│   ├── cache/                   # File-based cache (30min TTL)
│   ├── cli/                     # Arg parsing, help, setup wizard
│   ├── config/                  # Config struct, URL builder
│   ├── credentials/             # OS keyring (zalando/go-keyring)
│   ├── service/                 # Weather service (orchestrates cache + API)
│   ├── ui/                      # Icons, borders, ANSI colors
│   └── weather/                 # Display formatting
└── .github/workflows/           # CI (lint + test), Release (goreleaser)
```

## Branching
- **`dev`**: All work branches from `dev`
- **`main`**: Staging/production - never branch from or commit directly
- **Naming**: `feat/`, `fix/`, `refactor/`, `docs/`, `chore/`
- **PRs**: Target `dev`. Use PR stacks for dependent changes.

## Commands
```bash
go build -o weather-cli .        # Build
go run . [location]              # Run
go test ./...                    # Test all
go test -run TestName ./pkg      # Single test
go test -race ./...              # Race detection (CI default)
golangci-lint run                # Lint (v2)
```

## Code Map

| Symbol | Location | Role |
|--------|----------|------|
| `main` | main.go | CLI dispatch via switch |
| `Weather` | service/weather.go | Orchestrates cache + API fetch |
| `Client` | api/weather/client.go | HTTP client, URL building |
| `Cache` | cache/cache.go | File-based JSON cache |
| `Config` | config/config.go | API key, location, days |
| `Display` | weather/display.go | Formats weather output |
| `Parse` | cli/cli.go | Arg parsing, returns Command |

## Code Style
- **Imports**: stdlib > external > internal
- **Errors**: Wrap with context: `fmt.Errorf("action: %w", err)`
- **Types**: JSON tags on API structs; `float32` for weather numerics
- **Strings**: `strings.Builder` for concatenation

## Testing Instructions
- Run `go test -race ./...` before every commit
- Tests are colocated: `foo.go` + `foo_test.go`
- Use table-driven tests for multiple cases
- Mock via interfaces: `WeatherFetcher`, `WeatherCache` in service/weather.go
- Check `export_test.go` files - they expose internals for testing
- Fix any test failures before submitting PR

## Before Committing
1. Run `golangci-lint run` - fix all warnings
2. Run `go test -race ./...` - all tests must pass
3. Run `go build .` - ensure it compiles
4. Check `git diff` - no debug code, no hardcoded keys

## PR Instructions
- Branch from `dev`, target PR to `dev`
- Title format: `type(scope): description` (e.g., `fix(cache): handle corrupted files`)
- Run lint and tests before pushing - CI will fail otherwise
- Keep changes focused - split unrelated changes into separate PRs

## Security - Never Do This
- Never hardcode API keys or secrets
- Never commit `.env` files or credentials
- Never use `// nolint` without explanation
- Never suppress errors silently - always log or return

## Gotchas
- API key stored in OS keyring, not config file - use `credentials.Get()`
- Cache uses `float32` for weather data, not `float64`
- `internal/` packages can't be imported externally
- golangci-lint v2 syntax differs from v1 - check `.golangci.yml`

## CI
- **Trigger**: Push to main, PRs to main/dev
- **Jobs**: `lint` (golangci-lint v2.6), `test` (go test -race)
- **Release**: goreleaser on tag push

## API Key Setup
```bash
weather-cli --setup              # Interactive setup (stores in OS keyring)
WEATHER_API_KEY=xxx weather-cli  # Env var override
```
Get key: https://www.weatherapi.com/
