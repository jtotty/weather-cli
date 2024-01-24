package weather

type Data struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float32 `json:"temp_c"`
		FeelsLike float32 `json:"feelslike_c"`
		Humidity  float32 `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		AirQuality struct {
			PM25 float32 `json:"pm2_5"`
			PM10 float32 `json:"pm10"`
		} `json:"air_quality"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float32 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float32 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}
