package models

type Response struct {
	Status  string         `json:status`
	Payload interface{}    `json:payload`
	Errors  []ErrorPayload `json:errors`
}

type ErrorPayload struct {
	Detail string
}

type Temperatura struct {
	TemperaturaPromedio float64 `json:"temperaturaPromedio"`
	DesviacionEstandar  float64 `json:"desviacionEstandar"`
}
