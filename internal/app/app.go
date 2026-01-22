// Package app provides the application orchestration layer with dependency injection.
package app

import (
	"context"
	"io"

	"github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/config"
	weatherdisplay "github.com/jtotty/weather-cli/internal/weather"
)

// WeatherService defines the interface for getting weather data.
type WeatherService interface {
	GetWeather(ctx context.Context) (*weather.Response, error)
}

// App contains the application dependencies and orchestrates the weather workflow.
type App struct {
	config  *config.Config
	service WeatherService
	output  io.Writer
}

// NewApp creates an App with the provided dependencies.
// All dependencies are explicit - the caller is responsible for creating them.
func NewApp(cfg *config.Config, svc WeatherService, output io.Writer) *App {
	return &App{
		config:  cfg,
		service: svc,
		output:  output,
	}
}

// Run executes the weather fetch and display workflow.
func (a *App) Run(ctx context.Context, location string) error {
	if location != "" {
		a.config.SetLocation(location)
	}

	data, err := a.service.GetWeather(ctx)
	if err != nil {
		return err
	}

	display, err := weatherdisplay.NewDisplay(data, a.config.IsLocal)
	if err != nil {
		return err
	}

	display.RenderTo(a.output)
	return nil
}
