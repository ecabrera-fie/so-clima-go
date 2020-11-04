package providers

import (
	"context"
	"so2-clima/client"
	"so2-clima/models"
	"so2-clima/utils"
)

const (
	OFFLINE_STATUS = "OFFLINE"
	ONLINE_STATUS = "ONLINE"
	ERROR_STATUS = "ERROR"
	WARNING_STATUS = "WARNING"
	OK_STATUS = "OK"
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
	var providers []models.WeatherProvider

	response := models.Response{
		Status: OK_STATUS,
		Errors: []models.ErrorPayload{},
	}

	provider := models.WeatherProvider{ Name: p.accuweatherProvider.C.Name, Status: ONLINE_STATUS }
	accu, err := p.getLocationTemperatureFromAccuweatherAPI(geo)
	if err != nil {
		response.Status = WARNING_STATUS
		provider.Status = OFFLINE_STATUS
		response.Errors = append(response.Errors, models.ErrorPayload{Detail: err.Error()})
	} else {
		tempData = append(tempData, accu.Temperature.Metric.Value)
		sumTemp += accu.Temperature.Metric.Value
	}
	providers = append(providers, provider)

	provider = models.WeatherProvider{ Name: p.openweathermapProvider.C.Name, Status: ONLINE_STATUS }
	open, err := p.openweathermapProvider.GetOpenweatherMapWeather(p.Context, geo)
	if err != nil {
		response.Status = ERROR_STATUS
		provider.Status = OFFLINE_STATUS
		response.Errors = append(response.Errors, models.ErrorPayload{Detail: err.Error()})
	} else {
		tempData = append(tempData, open.Main.Temp)
		sumTemp += open.Main.Temp
	}
	providers = append(providers, provider)

	mintemp, _ := utils.Min(tempData)
	maxtemp, _ := utils.Max(tempData)

	if response.Status != ERROR_STATUS {
		weather := models.Weather{
			MeanTemp:          sumTemp / float64(len(tempData)),
			StandardDeviation: utils.StandardDeviation(tempData),
			WeatherProviders: providers,
			MaxTemp: maxtemp,
			MinTemp: mintemp,
		}

		response.Payload = weather
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
