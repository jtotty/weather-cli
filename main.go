package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/jtotty/weather-cli/internal/cli"
	"github.com/jtotty/weather-cli/internal/credentials"
	"github.com/jtotty/weather-cli/internal/service"
	"github.com/jtotty/weather-cli/internal/weather"
)

var version = "dev"

func main() {
	cmd := cli.Parse(os.Args)
	store := credentials.NewEnvOverride(credentials.NewKeyring())

	switch cmd.Type {
	case cli.CommandHelp:
		cli.PrintHelp(version)
	case cli.CommandVersion:
		cli.PrintVersion(version)
	case cli.CommandSetup:
		if err := cli.RunSetup(store, cli.ReadKeyFromTerminal); err != nil {
			cli.ExitWithError(fmt.Errorf("setup failed: %w", err))
		}
	case cli.CommandDeleteKey:
		if err := cli.RunDeleteKey(store); err != nil {
			cli.ExitWithError(fmt.Errorf("failed to delete API key: %w", err))
		}
	case cli.CommandWeather:
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()
		runWeather(ctx, store, cmd.Location)
	}
}

func runWeather(ctx context.Context, store credentials.Store, location string) {
	cfg, err := cli.LoadConfig(store, cli.ReadKeyFromTerminal)
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
			fmt.Fprintln(os.Stderr, "\nRequest canceled.")
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
