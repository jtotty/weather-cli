package cli

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantType     CommandType
		wantLocation string
	}{
		{
			name:         "no arguments returns weather command",
			args:         []string{"weather-cli"},
			wantType:     CommandWeather,
			wantLocation: "",
		},
		{
			name:         "location argument",
			args:         []string{"weather-cli", "London"},
			wantType:     CommandWeather,
			wantLocation: "London",
		},
		{
			name:         "location with spaces",
			args:         []string{"weather-cli", "New York"},
			wantType:     CommandWeather,
			wantLocation: "New York",
		},
		{
			name:         "zip code location",
			args:         []string{"weather-cli", "10001"},
			wantType:     CommandWeather,
			wantLocation: "10001",
		},
		{
			name:         "coordinates location",
			args:         []string{"weather-cli", "51.5,-0.1"},
			wantType:     CommandWeather,
			wantLocation: "51.5,-0.1",
		},
		{
			name:     "--help flag",
			args:     []string{"weather-cli", "--help"},
			wantType: CommandHelp,
		},
		{
			name:     "-h flag",
			args:     []string{"weather-cli", "-h"},
			wantType: CommandHelp,
		},
		{
			name:     "--version flag",
			args:     []string{"weather-cli", "--version"},
			wantType: CommandVersion,
		},
		{
			name:     "-v flag",
			args:     []string{"weather-cli", "-v"},
			wantType: CommandVersion,
		},
		{
			name:     "--setup flag",
			args:     []string{"weather-cli", "--setup"},
			wantType: CommandSetup,
		},
		{
			name:     "--delete-key flag",
			args:     []string{"weather-cli", "--delete-key"},
			wantType: CommandDeleteKey,
		},
		{
			name:     "unknown flag shows help",
			args:     []string{"weather-cli", "--unknown"},
			wantType: CommandHelp,
		},
		{
			name:     "single dash unknown flag shows help",
			args:     []string{"weather-cli", "-x"},
			wantType: CommandHelp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.args)

			if got.Type != tt.wantType {
				t.Errorf("Parse() Type = %v, want %v", got.Type, tt.wantType)
			}

			if got.Location != tt.wantLocation {
				t.Errorf("Parse() Location = %q, want %q", got.Location, tt.wantLocation)
			}
		})
	}
}

func TestCommandType_String(t *testing.T) {
	// Verify command type constants have expected values
	if CommandWeather != 0 {
		t.Errorf("CommandWeather = %d, want 0", CommandWeather)
	}
	if CommandHelp != 1 {
		t.Errorf("CommandHelp = %d, want 1", CommandHelp)
	}
	if CommandVersion != 2 {
		t.Errorf("CommandVersion = %d, want 2", CommandVersion)
	}
	if CommandSetup != 3 {
		t.Errorf("CommandSetup = %d, want 3", CommandSetup)
	}
	if CommandDeleteKey != 4 {
		t.Errorf("CommandDeleteKey = %d, want 4", CommandDeleteKey)
	}
}
