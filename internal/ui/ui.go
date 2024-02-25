package ui

import (
	"fmt"
	"strings"

	"github.com/enescakir/emoji"
)

var Icons = map[string]string{
	"wind":     emoji.LeafFlutteringInWind.String(),
	"humidity": emoji.Droplet.String(),
	"sunrise":  emoji.Sunrise.String(),
	"sunset":   emoji.Sunset.String(),
}

var WeatherIcons = map[string]string{
	"clear":                                    emoji.NightWithStars.String(),
	"sunny":                                    emoji.Sun.String(),
	"partly_cloudy":                            emoji.SunBehindCloud.String(),
	"cloudy":                                   "",
	"overcast":                                 "",
	"mist":                                     emoji.Fog.String(),
	"patchy_rain_possible":                     "",
	"patchy_snow_possible":                     "",
	"patchy_sleet_possible":                    "",
	"patchy_freezing_drizzle_possible":         "",
	"thundery_outbreaks_possible":              "",
	"blowing_snow":                             "",
	"blizzard":                                 "",
	"fog":                                      "",
	"freezing_fog":                             "",
	"patchy_light_drizzle":                     "",
	"light_drizzle":                            "",
	"freezing_drizzle":                         "",
	"heavy_freezing_drizzle":                   "",
	"patchy_light_rain":                        "",
	"light_rain":                               "",
	"moderate_rain_at_times":                   "",
	"moderate_rain":                            "",
	"heavy_rain_at_times":                      "",
	"heavy_rain":                               "",
	"light_freezing_rain":                      "",
	"moderate_or_heavy_freezing_rain":          "",
	"light_sleet":                              "",
	"moderate_or_heavy_sleet":                  "",
	"patchy_light_snow":                        "",
	"light_snow":                               "",
	"patchy_moderate_snow":                     "",
	"moderate_snow":                            "",
	"patchy_heavy_snow":                        "",
	"heavy_snow":                               "",
	"ice_pellets":                              "",
	"light_rain_shower":                        "",
	"moderate_or_heavy_rain_shower":            "",
	"torrential_rain_shower":                   "",
	"light_sleet_showers":                      "",
	"moderate_or_heavy_sleet_showers":          "",
	"light_snow_showers":                       "",
	"moderate_or_heavy_snow_showers":           "",
	"light_showers_of_ice_pellets":             "",
	"moderate_or_heavy_showers_of_ice_pellets": "",
	"patchy_light_rain_with_thunder":           "",
	"moderate_or_heavy_rain_with_thunder":      "",
	"patchy_light_snow_with_thunder":           "",
	"moderate_or_heavy_snow_with_thunder":      "",
}

var aqiIcons = map[string]string{
	"good":           emoji.GreenCircle.String(),
	"moderate":       emoji.YellowCircle.String(),
	"sensitive":      emoji.OrangeCircle.String(),
	"unhealthy":      emoji.RedCircle.String(),
	"very_unhealthy": emoji.PurpleCircle.String(),
	"hazardous":      emoji.Skull.String(),
}

func CreateBorder(maxLen int) string {
	border := strings.Builder{}
	for i := 0; i < maxLen; i++ {
		border.WriteString("-")
	}

	return border.String()
}

func Spacer() {
	fmt.Print("\n\n")
}

func GetWeatherIcon(name string) string {
	key := strings.ReplaceAll(strings.ToLower(name), " ", "_")

    if val, ok := WeatherIcons[key]; ok {
        return val
    }

    return ""
}

func GetAqiIcon(val float32) string {
	aqi := int(val)

	switch {
	case aqi >= 0 && aqi <= 50:
		return aqiIcons["good"]
	case aqi >= 51 && aqi <= 100:
		return aqiIcons["moderate"]
	case aqi >= 101 && aqi <= 150:
		return aqiIcons["sensitive"]
	case aqi >= 151 && aqi <= 200:
		return aqiIcons["unhealthy"]
	case aqi >= 201 && aqi <= 300:
		return aqiIcons["very_unhealthy"]
	case aqi > 300:
		return aqiIcons["hazardous"]
	default:
		return "" // For negative AQI values
	}
}
