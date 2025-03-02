package models

type Weather struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
	Alerts   Alerts   `json:"alerts"`
}

type Location struct {
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	TzID      string  `json:"tz_id"`
	LocalTime string  `json:"localtime"`
}

type Current struct {
	TempC         float32    `json:"temp_c"`
	TempF         float32    `json:"temp_f"`
	FeelsLike     float32    `json:"feelslike_c"`
	FeelsLikeF    float32    `json:"feelslike_f"`
	Humidity      float32    `json:"humidity"`
	WindSpeed     float32    `json:"wind_mph"`
	WindSpeedKph  float32    `json:"wind_kph"`
	WindDirection string     `json:"wind_dir"`
	Condition     Condition  `json:"condition"`
	AirQuality    AirQuality `json:"air_quality"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type AirQuality struct {
	CO    float32 `json:"co"`
	NO2   float32 `json:"no2"`
	O3    float32 `json:"o3"`
	SO2   float32 `json:"so2"`
	PM25  float32 `json:"pm2_5"`
	PM10  float32 `json:"pm10"`
	Index int     `json:"us-epa-index"`
}

type Forecast struct {
	Forecastday []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date       string     `json:"date"`
	DateEpoch  int64      `json:"date_epoch"`
	Day        Day        `json:"day"`
	Astro      Astro      `json:"astro"`
	Hour       []Hour     `json:"hour"`
	AirQuality AirQuality `json:"air_quality"`
}

type Day struct {
	MaxTempC     float32   `json:"maxtemp_c"`
	MinTempC     float32   `json:"mintemp_c"`
	AvgTempC     float32   `json:"avgtemp_c"`
	Condition    Condition `json:"condition"`
	ChanceOfRain float32   `json:"daily_chance_of_rain"`
}

type Hour struct {
	TimeEpoch    int64     `json:"time_epoch"`
	Time         string    `json:"time"`
	TempC        float32   `json:"temp_c"`
	Condition    Condition `json:"condition"`
	ChanceOfRain float32   `json:"chance_of_rain"`
	ChanceOfSnow float32   `json:"chance_of_snow"`
}

type Astro struct {
	Sunrise   string `json:"sunrise"`
	Sunset    string `json:"sunset"`
	Moonrise  string `json:"moonrise"`
	Moonset   string `json:"moonset"`
	MoonPhase string `json:"moon_phase"`
}

type Alerts struct {
	Alert []Alert `json:"alert"`
}

type Alert struct {
	Headline    string `json:"headline"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Event       string `json:"event"`
	Description string `json:"desc"`
}
