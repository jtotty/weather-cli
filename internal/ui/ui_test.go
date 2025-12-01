package ui

import (
	"fmt"
	"testing"
)

func TestGetAqiIcon_Thresholds(t *testing.T) {
	// Test that icons change at each AQI threshold boundary
	// We test boundary values to ensure correct categorization
	tests := []struct {
		aqi      float32
		category string
	}{
		{0, "good"},
		{50, "good"},
		{51, "moderate"},
		{100, "moderate"},
		{101, "sensitive"},
		{150, "sensitive"},
		{151, "unhealthy"},
		{200, "unhealthy"},
		{201, "very_unhealthy"},
		{300, "very_unhealthy"},
		{301, "hazardous"},
	}

	// Track previous icon to verify it changes at boundaries
	var prevIcon string
	var prevCategory string

	for _, tt := range tests {
		t.Run(fmt.Sprintf("AQI_%.0f_%s", tt.aqi, tt.category), func(t *testing.T) {
			icon := GetAqiIcon(tt.aqi)

			if icon == "" {
				t.Errorf("GetAqiIcon(%.0f) returned empty string", tt.aqi)
			}

			// Verify icon changes when category changes
			if prevCategory != "" && prevCategory != tt.category && icon == prevIcon {
				t.Errorf("Icon should change between %s and %s categories", prevCategory, tt.category)
			}

			prevIcon = icon
			prevCategory = tt.category
		})
	}
}

func TestGetAqiIcon_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		aqi      float32
		notEmpty bool
	}{
		{"zero", 0, true},
		{"negative", -10, true},
		{"very high", 999, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			icon := GetAqiIcon(tt.aqi)

			if tt.notEmpty && icon == "" {
				t.Errorf("GetAqiIcon(%.0f) returned empty string, expected non-empty", tt.aqi)
			}
		})
	}
}

func TestCreateBorder_Length(t *testing.T) {
	tests := []struct {
		length int
		want   int
	}{
		{0, 0},
		{1, 1},
		{5, 5},
		{100, 100},
		{-1, 0}, // Negative should return empty string (length 0)
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("len_%d", tt.length), func(t *testing.T) {
			border := CreateBorder(tt.length)

			if len(border) != tt.want {
				t.Errorf("CreateBorder(%d) length = %d, want %d", tt.length, len(border), tt.want)
			}

			// Verify it's all dashes
			for _, c := range border {
				if c != '-' {
					t.Errorf("CreateBorder(%d) contains non-dash character: %c", tt.length, c)
				}
			}
		})
	}
}
