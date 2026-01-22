package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/jtotty/weather-cli/internal/app"
	"github.com/jtotty/weather-cli/internal/cli"
	"github.com/jtotty/weather-cli/internal/config"
	"github.com/jtotty/weather-cli/internal/credentials"
	"github.com/jtotty/weather-cli/internal/service"
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

	cache, fetcher := service.NewDefaultDeps(cfg)
	svc := service.NewWeather(cfg, cache, fetcher)
	application := app.NewApp(cfg, svc, os.Stdout)
	if err := application.Run(ctx, location); err != nil {
		if errors.Is(err, context.Canceled) {
			fmt.Fprintln(os.Stderr, "\nRequest canceled.")
			os.Exit(130)
		}
		cli.ExitWithError(fmt.Errorf("error fetching weather: %w", err))
	}
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
