package ui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/enescakir/emoji"
)

const (
	AQIGood          = 50
	AQIModerate      = 100
	AQISensitive     = 150
	AQIUnhealthy     = 200
	AQIVeryUnhealthy = 300
)

var icons = map[string]emoji.Emoji{
	"wind":     emoji.LeafFlutteringInWind,
	"humidity": emoji.Droplet,
	"sunrise":  emoji.Sunrise,
	"sunset":   emoji.Sunset,
}

var weatherIcons = map[string]emoji.Emoji{
	"clear":                                    emoji.NightWithStars,
	"sunny":                                    emoji.Sun,
	"partly_cloudy":                            emoji.SunBehindCloud,
	"cloudy":                                   emoji.Cloud,
	"overcast":                                 emoji.Cloud,
	"mist":                                     emoji.Fog,
	"patchy_rain_possible":                     emoji.CloudWithRain,
	"patchy_rain_nearby":                       emoji.CloudWithRain,
	"patchy_snow_possible":                     emoji.CloudWithSnow,
	"patchy_sleet_possible":                    emoji.CloudWithRain + emoji.Snowflake,
	"patchy_freezing_drizzle_possible":         emoji.CloudWithRain + emoji.Ice,
	"thundery_outbreaks_possible":              emoji.CloudWithLightning,
	"blowing_snow":                             emoji.CloudWithSnow,
	"blizzard":                                 emoji.CloudWithSnow + emoji.DashingAway,
	"fog":                                      emoji.Fog,
	"freezing_fog":                             emoji.Fog + emoji.Ice,
	"patchy_light_drizzle":                     emoji.CloudWithRain,
	"light_drizzle":                            emoji.CloudWithRain,
	"freezing_drizzle":                         emoji.CloudWithRain + emoji.Ice,
	"heavy_freezing_drizzle":                   emoji.CloudWithRain + emoji.Ice,
	"patchy_light_rain":                        emoji.CloudWithRain,
	"light_rain":                               emoji.CloudWithRain,
	"moderate_rain_at_times":                   emoji.CloudWithRain,
	"moderate_rain":                            emoji.CloudWithRain,
	"heavy_rain_at_times":                      emoji.CloudWithRain,
	"heavy_rain":                               emoji.CloudWithRain,
	"light_freezing_rain":                      emoji.CloudWithRain + emoji.Ice,
	"moderate_or_heavy_freezing_rain":          emoji.CloudWithRain + emoji.Ice,
	"light_sleet":                              emoji.CloudWithRain,
	"moderate_or_heavy_sleet":                  emoji.CloudWithRain,
	"patchy_light_snow":                        emoji.CloudWithSnow,
	"light_snow":                               emoji.CloudWithSnow,
	"patchy_moderate_snow":                     emoji.CloudWithSnow,
	"moderate_snow":                            emoji.CloudWithSnow,
	"patchy_heavy_snow":                        emoji.CloudWithSnow,
	"heavy_snow":                               emoji.CloudWithSnow,
	"ice_pellets":                              emoji.Ice,
	"light_rain_shower":                        emoji.CloudWithRain,
	"moderate_or_heavy_rain_shower":            emoji.CloudWithRain,
	"torrential_rain_shower":                   emoji.CloudWithRain,
	"light_sleet_showers":                      emoji.CloudWithRain + emoji.Ice,
	"moderate_or_heavy_sleet_showers":          emoji.CloudWithRain + emoji.Ice,
	"light_snow_showers":                       emoji.CloudWithSnow,
	"moderate_or_heavy_snow_showers":           emoji.CloudWithSnow,
	"light_showers_of_ice_pellets":             emoji.CloudWithSnow + emoji.Ice,
	"moderate_or_heavy_showers_of_ice_pellets": emoji.CloudWithSnow + emoji.Ice,
	"patchy_light_rain_with_thunder":           emoji.CloudWithLightningAndRain,
	"moderate_or_heavy_rain_with_thunder":      emoji.CloudWithLightningAndRain,
	"patchy_light_snow_with_thunder":           emoji.CloudWithLightning + emoji.Snowflake,
	"moderate_or_heavy_snow_with_thunder":      emoji.CloudWithLightning + emoji.Snowflake,
}

var aqiIcons = map[string]emoji.Emoji{
	"good":           emoji.GreenCircle,
	"moderate":       emoji.YellowCircle,
	"sensitive":      emoji.OrangeCircle,
	"unhealthy":      emoji.RedCircle,
	"very_unhealthy": emoji.PurpleCircle,
	"hazardous":      emoji.Skull,
}

func GetIcon(name string) string {
	key := strings.ToLower(name)

	if icon, ok := icons[key]; ok {
		return icon.String()
	}

	return ""
}

func GetWeatherIcon(name string) string {
	key := strings.TrimSpace(strings.ToLower(name))
	key = strings.ReplaceAll(key, " ", "_")

	if icon, ok := weatherIcons[key]; ok {
		return icon.String()
	}

	return "Err: Icon not loaded"
}

func GetAqiIcon(num float32) string {
	aqi := int(num)
	if aqi < 0 {
		return emoji.QuestionMark.String()
	}

	var icon emoji.Emoji

	switch {
	case aqi <= AQIGood:
		icon = aqiIcons["good"]
	case aqi <= AQIModerate:
		icon = aqiIcons["moderate"]
	case aqi <= AQISensitive:
		icon = aqiIcons["sensitive"]
	case aqi <= AQIUnhealthy:
		icon = aqiIcons["unhealthy"]
	case aqi <= AQIVeryUnhealthy:
		icon = aqiIcons["very_unhealthy"]
	default:
		icon = aqiIcons["hazardous"]
	}

	return icon.String()
}

func CreateBorder(maxLen int) string {
	border := strings.Builder{}
	for i := 0; i < maxLen; i++ {
		border.WriteString("-")
	}

	return border.String()
}

func Spacer() {
	fmt.Fprint(os.Stdout, "\n\n")
}

func SpacerTo(w io.Writer) {
	fmt.Fprint(w, "\n\n")
}
