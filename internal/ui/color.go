package ui

import "fmt"

// ANSI color codes (24-bit true color foreground)
// Format: \033[38;2;R;G;Bm
const ColorReset = "\033[0m"

// Temperature color scale (Fahrenheit ranges with hex converted to RGB)
// Gradient from cold (blues) to hot (reds)
var tempColors = []struct {
	maxTempF float32
	color    string
}{
	{-100, "\033[38;2;228;240;255m"},
	{-60, "\033[38;2;228;240;255m"},
	{-55, "\033[38;2;219;233;251m"},
	{-50, "\033[38;2;211;226;247m"},
	{-45, "\033[38;2;203;220;244m"},
	{-40, "\033[38;2;192;213;237m"},
	{-35, "\033[38;2;184;206;232m"},
	{-30, "\033[38;2;176;199;231m"},
	{-25, "\033[38;2;167;192;227m"},
	{-20, "\033[38;2;157;184;222m"},
	{-15, "\033[38;2;146;175;213m"},
	{-10, "\033[38;2;136;165;206m"},
	{-5, "\033[38;2;128;155;195m"},
	{0, "\033[38;2;118;145;185m"},
	{5, "\033[38;2;96;124;167m"},
	{10, "\033[38;2;86;114;156m"},
	{15, "\033[38;2;77;102;145m"},
	{20, "\033[38;2;65;93;135m"},
	{25, "\033[38;2;57;82;127m"},
	{30, "\033[38;2;47;72;117m"},
	{35, "\033[38;2;39;67;111m"},
	{40, "\033[38;2;36;79;120m"},
	{45, "\033[38;2;39;92;128m"},
	{50, "\033[38;2;39;103;138m"},
	{55, "\033[38;2;39;117;147m"},
	{60, "\033[38;2;68;128;144m"},
	{65, "\033[38;2;100;141;137m"},
	{70, "\033[38;2;135;155;132m"},
	{75, "\033[38;2;172;168;125m"},
	{80, "\033[38;2;195;171;117m"},
	{85, "\033[38;2;191;159;104m"},
	{90, "\033[38;2;195;139;83m"},
	{95, "\033[38;2;193;111;74m"},
	{100, "\033[38;2;175;77;78m"},
	{105, "\033[38;2;159;41;76m"},
	{110, "\033[38;2;135;32;62m"},
	{115, "\033[38;2;110;21;50m"},
	{120, "\033[38;2;87;11;37m"},
	{150, "\033[38;2;61;2;22m"},
}

// celsiusToFahrenheit converts Celsius to Fahrenheit
func celsiusToFahrenheit(c float32) float32 {
	return (c * 9 / 5) + 32
}

// ColorizeTemp returns a temperature string with ANSI color coding.
func ColorizeTemp(temp float32) string {
	color := getTempColor(temp)
	return fmt.Sprintf("%s%3.0f°C%s", color, temp, ColorReset)
}

// getTempColor returns the appropriate ANSI color code for a temperature in Celsius.
func getTempColor(temp float32) string {
	tempF := celsiusToFahrenheit(temp)

	for _, tc := range tempColors {
		if tempF <= tc.maxTempF {
			return tc.color
		}
	}
	// Above max, use the hottest color
	return tempColors[len(tempColors)-1].color
}
