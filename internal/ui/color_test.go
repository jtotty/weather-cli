package ui

import (
	"fmt"
	"strings"
	"testing"
)

func TestColorizeTemp_ReturnsColoredString(t *testing.T) {
	tests := []struct {
		name string
		temp float32
	}{
		{"extreme cold", -50},
		{"very cold", -30},
		{"cold", -10},
		{"freezing", 0},
		{"cool", 10},
		{"mild", 15},
		{"pleasant", 20},
		{"warm", 25},
		{"hot", 30},
		{"very hot", 40},
		{"extreme heat", 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorizeTemp(tt.temp)

			// Should have color prefix
			if !strings.HasPrefix(result, "\033[38;2;") {
				t.Errorf("ColorizeTemp(%v) = %q, want ANSI color prefix", tt.temp, result)
			}

			// Should have reset suffix
			if !strings.HasSuffix(result, ColorReset) {
				t.Errorf("ColorizeTemp(%v) = %q, want suffix %q", tt.temp, result, ColorReset)
			}

			// Should contain degree symbol
			if !strings.Contains(result, "°C") {
				t.Errorf("ColorizeTemp(%v) = %q, want to contain '°C'", tt.temp, result)
			}
		})
	}
}

func TestColorizeTemp_Format(t *testing.T) {
	tests := []struct {
		name     string
		temp     float32
		wantTemp string
	}{
		{"single digit", 5, "  5°C"},
		{"double digit", 25, " 25°C"},
		{"negative", -5, " -5°C"},
		{"rounds down", 15.4, " 15°C"},
		{"rounds up", 15.6, " 16°C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColorizeTemp(tt.temp)

			// Strip color codes to check formatted temp
			stripped := strings.TrimPrefix(result, getTempColor(tt.temp))
			stripped = strings.TrimSuffix(stripped, ColorReset)

			if stripped != tt.wantTemp {
				t.Errorf("ColorizeTemp(%v) formatted as %q, want %q", tt.temp, stripped, tt.wantTemp)
			}
		})
	}
}

func TestGetTempColor_GradientOrder(t *testing.T) {
	// Verify that colder temps get different colors than warmer temps
	coldColor := getTempColor(-20)
	coolColor := getTempColor(5)
	warmColor := getTempColor(25)
	hotColor := getTempColor(40)

	// All should be different colors
	colors := []string{coldColor, coolColor, warmColor, hotColor}
	seen := make(map[string]bool)
	for _, c := range colors {
		if seen[c] {
			t.Error("Expected different colors for different temperature ranges")
		}
		seen[c] = true
	}
}

func TestGetTempColor_ExtremeValues(t *testing.T) {
	// Test extreme temperatures return valid colors
	belowMin := getTempColor(-80)
	aboveMax := getTempColor(70)

	// Should return valid color codes
	if !strings.HasPrefix(belowMin, "\033[38;2;") {
		t.Errorf("getTempColor(-80) = %q, want ANSI color code", belowMin)
	}
	if !strings.HasPrefix(aboveMax, "\033[38;2;") {
		t.Errorf("getTempColor(70) = %q, want ANSI color code", aboveMax)
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		celsius    float32
		fahrenheit float32
	}{
		{0, 32},
		{100, 212},
		{-40, -40},
		{20, 68},
		{-20, -4},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v°C", tt.celsius), func(t *testing.T) {
			result := celsiusToFahrenheit(tt.celsius)
			if result != tt.fahrenheit {
				t.Errorf("celsiusToFahrenheit(%v) = %v, want %v", tt.celsius, result, tt.fahrenheit)
			}
		})
	}
}
