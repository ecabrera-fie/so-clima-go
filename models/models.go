package models

type Response struct {
	Status  string         `json:"status"`
	Payload interface{}    `json:"payload"`
	Errors  []ErrorPayload `json:"errors"`
}

type ErrorPayload struct {
	Detail string
}

type Weather struct {
	MeanTemp          float64           `json:"temperaturaPromedio"`
	StandardDeviation float64           `json:"desviacionEstandar"`
	MinTemp           float64           `json:"temperaturaMinima"`
	MaxTemp           float64           `json:"temperaturaMaxima"`
	WeatherProviders  []WeatherProvider `json:"proveedores"`
}

type WeatherProvider struct {
	Name   string `json:"nombre"`
	Status string `json:"estado"`
}