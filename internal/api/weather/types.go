package weather

type Response struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
	Alerts   Alerts   `json:"alerts"`
}

type Location struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	LocalTime string `json:"localtime"`
}

type Current struct {
	TempC         float32    `json:"temp_c"`
	FeelsLike     float32    `json:"feelslike_c"`
	Humidity      float32    `json:"humidity"`
	WindSpeed     float32    `json:"wind_mph"`
	WindDirection string     `json:"wind_dir"`
	Condition     Condition  `json:"condition"`
	AirQuality    AirQuality `json:"air_quality"`
}

type Condition struct {
	Text string `json:"text"`
}

type AirQuality struct {
	PM25 float32 `json:"pm2_5"`
	PM10 float32 `json:"pm10"`
}

type Forecast struct {
	Forecastday []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Hour       []Hour     `json:"hour"`
	AirQuality AirQuality `json:"air_quality"`
	Astro      Astro      `json:"astro"`
}

type Hour struct {
	TimeEpoch    int64     `json:"time_epoch"`
	TempC        float32   `json:"temp_c"`
	Condition    Condition `json:"condition"`
	ChanceOfRain float32   `json:"chance_of_rain"`
}

type Astro struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

type Alerts struct {
	Alert []Alert `json:"alert"`
}

type Alert struct {
	Event string `json:"event"`
	Desc  string `json:"desc"`
}
