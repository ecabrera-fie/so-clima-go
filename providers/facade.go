package providers

import (
	"context"
	"so2-clima/client"
	"so2-clima/models"
	"so2-clima/utils"
)

type ProvidersFacade struct {
	Context                context.Context
	accuweatherProvider    *AccuweatherClientWrapper
	openweathermapProvider *OpenweatherClientWrapper
}

func NewDistributedWeatherProvider(ctx context.Context) *ProvidersFacade {
	return &ProvidersFacade{
		Context:                ctx,
		accuweatherProvider:    NewAccuweatherClient(),
		openweathermapProvider: NewOpenweatherMapClient(),
	}
}

func (p *ProvidersFacade) GetTemperatureDataByGeolocation(geo *client.Geoposition) models.Response {
	var sumTemp float64
	var tempData []float64

	response := models.Response{
		Status: "OK",
		Errors: []models.ErrorPayload{},
	}

	accu, err := p.getLocationTemperatureFromAccuweatherAPI(geo)
	if err != nil {
		response.Status = "WARNING"
		response.Errors = append(response.Errors, models.ErrorPayload{Detail: err.Error()})
	} else {
		tempData = append(tempData, accu.Temperature.Metric.Value)
		sumTemp += accu.Temperature.Metric.Value
	}
	open, err := p.openweathermapProvider.GetOpenweatherMapWeather(p.Context, geo)
	if err != nil {
		response.Status = "ERROR"
		response.Errors = append(response.Errors, models.ErrorPayload{Detail: err.Error()})
	} else {
		tempData = append(tempData, open.Main.Temp)
		sumTemp += open.Main.Temp
	}

	if response.Status != "ERROR" {
		temp := models.Temperatura{
			TemperaturaPromedio: sumTemp / float64(len(tempData)),
			DesviacionEstandar:  utils.StandardDeviation(tempData),
		}

		response.Payload = temp
	}

	return response
}

func (p *ProvidersFacade) getLocationTemperatureFromAccuweatherAPI(geo *client.Geoposition) (*AccuweatherCurrentWeather, error) {
	city, err := p.accuweatherProvider.GetAccuweatherCityByGeoposition(p.Context, geo)
	if err != nil {
		return nil, err
	}

	clima, err := p.accuweatherProvider.GetAccuweatherCurrentWeatherByCityKey(p.Context, city.Key)
	if err != nil {
		return nil, err
	}

	return clima, nil
}
