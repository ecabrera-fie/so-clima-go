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

func (p *ProvidersFacade) GetTemperatureDataByGeolocation(geo *client.Geoposition) (models.Temperatura, error) {
	accu, _ := p.getLocationTemperatureFromAccuweatherAPI(geo)
	open, _ := p.openweathermapProvider.GetOpenweatherMapWeather(p.Context, geo)

	temp := models.Temperatura{
		TemperaturaPromedio: (accu.Temperature.Metric.Value + open.Main.Temp) / 2,
		DesviacionEstandar:  utils.StandardDeviation(accu.Temperature.Metric.Value, open.Main.Temp),
	}

	return temp, nil
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
