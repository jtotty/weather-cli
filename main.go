package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/jtotty/weather-cli/internal/cli"
	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/credentials"
	"github.com/jtotty/weather-cli/internal/service"
	"github.com/jtotty/weather-cli/internal/weather"
)

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
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		runWeather(ctx, cmd.Location)
	}
}

func runWeather(ctx context.Context, location string) {
	cfg, err := loadConfig()
	if err != nil {
		cli.ExitWithError(err)
	}

	if location != "" {
		cfg.SetLocation(location)
	}

	svc := service.NewWeather(cfg)
	data, err := svc.GetWeather(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, "\nRequest cancelled.")
			os.Exit(130)
		}
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
