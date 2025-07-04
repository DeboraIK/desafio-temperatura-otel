package dto

type Temperature struct {
	City   string  `json:"city"`
	Temp_C float32 `json:"temp_c"`
	Temp_F float32 `json:"temp_f"`
	Temp_K float32 `json:"temp_k"`
}

type OpenMeteoResponse struct {
	CurrentWeather struct {
		Temperature float32 `json:"temperature"`
	} `json:"current_weather"`
}
