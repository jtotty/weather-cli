# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Quick Reference

See `AGENTS.md` for the full project reference including architecture, code map, and conventions.

## Commands

```bash
go build -o weather-cli .        # Build
go run . [location]              # Run
go test ./...                    # Test all
go test -run TestName ./pkg      # Single test
go test -race ./...              # Test with race detection (CI default)
golangci-lint run                # Lint (v2 syntax)
```

## Before Committing

1. `golangci-lint run` - fix all warnings
2. `go test -race ./...` - all tests must pass
3. `go build .` - ensure it compiles

## Architecture

CLI weather tool using WeatherAPI.com. Go 1.25, OS keyring for credentials, file-based cache.

**Data flow**: `main.go` → `cli.Parse` → `service.Weather` → (cache hit? return : `api/weather.Client.Fetch`) → `weather.Display.Render`

**Key interfaces for testing** (in `service/weather.go`):
- `WeatherFetcher` - mock API calls
- `WeatherCache` - mock cache

## Branching

- Branch from `dev`, target PRs to `dev`
- Never commit directly to `main`
- Prefix: `feat/`, `fix/`, `refactor/`, `docs/`, `chore/`

## Agent skills

### Issue tracker

Issues and specs live as local markdown files under `.scratch/<feature>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Default label vocabulary (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.
