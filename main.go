package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jtotty/weather-cli/internal/cli"
	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/credentials"
	"github.com/jtotty/weather-cli/internal/service"
	"github.com/jtotty/weather-cli/internal/weather"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	cmd := cli.Parse(os.Args)

	switch cmd.Type {
	case cli.CommandHelp:
		cli.PrintHelp(version)
	case cli.CommandVersion:
		cli.PrintVersion(version)
	case cli.CommandSetup:
		if err := cli.RunSetup(); err != nil {
			cli.ExitWithError(fmt.Errorf("setup failed: %w", err))
		}
	case cli.CommandDeleteKey:
		if err := cli.RunDeleteKey(); err != nil {
			cli.ExitWithError(fmt.Errorf("failed to delete API key: %w", err))
		}
	case cli.CommandWeather:
		runWeather(cmd.Location)
	}
}

func runWeather(location string) {
	cfg, err := loadConfig()
	if err != nil {
		cli.ExitWithError(err)
	}

	if location != "" {
		cfg.SetLocation(location)
	}

	svc := service.NewWeather(cfg)
	data, err := svc.GetWeather()
	if err != nil {
		cli.ExitWithError(fmt.Errorf("error fetching weather: %w", err))
	}

	display, err := weather.NewDisplay(data, cfg.IsLocal)
	if err != nil {
		cli.ExitWithError(fmt.Errorf("error creating display: %w", err))
	}

	display.Render()
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err == nil {
		return cfg, nil
	}

	if errors.Is(err, credentials.ErrNoAPIKey) {
		fmt.Println("No API key configured.")
		fmt.Println()
		if setupErr := cli.RunSetup(); setupErr != nil {
			return nil, fmt.Errorf("setup failed: %w", setupErr)
		}
		return config.New()
	}

	return nil, fmt.Errorf("error loading config: %w", err)
}
