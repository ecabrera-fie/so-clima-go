package providers

import (
	"context"
	"sync"

	"../client"
	"../models"
	"../utils"
)

var waitGroup sync.WaitGroup

const (
	OFFLINE_STATUS = "OFFLINE"
	ONLINE_STATUS  = "ONLINE"
	ERROR_STATUS   = "ERROR"
	WARNING_STATUS = "WARNING"
	OK_STATUS      = "OK"
)

type standardWeatherInfo struct {
	temperature float64
}

type ProvidersFacade struct {
	Context                context.Context
	accuweatherProvider    *AccuweatherClientWrapper
	openweathermapProvider *OpenweatherClientWrapper
	climacellProvider      *ClimacellClientWrapper
}

func NewDistributedWeatherProvider(ctx context.Context) *ProvidersFacade {
	return &ProvidersFacade{
		Context:                ctx,
		accuweatherProvider:    NewAccuweatherClient(),
		openweathermapProvider: NewOpenweatherMapClient(),
		climacellProvider:      NewClimacellClient(),
	}
}

func (p *ProvidersFacade) GetTemperatureDataByGeolocation(geo *client.Geoposition) models.Response {
	var sumTemp float64
	var tempData []float64
	var providers []models.WeatherProvider
	channel := make(chan standardWeatherInfo)

	response := models.Response{}

	waitGroup.Add(3)

	go func() {
		waitGroup.Wait()
		close(channel)
	}()

	accuProvider := models.WeatherProvider{Name: p.accuweatherProvider.C.Name, Status: ONLINE_STATUS}
	go getLocationTemperatureFromAccuweatherAPI(p.Context, geo, &accuProvider, channel, p.accuweatherProvider)

	openProvider := models.WeatherProvider{Name: p.openweathermapProvider.C.Name, Status: ONLINE_STATUS}
	go getWeatherFromOpenWeatherAPI(p.Context, geo, &openProvider, channel, p.openweathermapProvider)

	climacellProvider := models.WeatherProvider{Name: p.climacellProvider.C.Name, Status: ONLINE_STATUS}
	go getWeatherFromClimacellAPI(p.Context, geo, &climacellProvider, channel, p.climacellProvider)

	for weatherInfo := range channel {
		tempData = append(tempData, weatherInfo.temperature)
		sumTemp += weatherInfo.temperature
	}

	providers = append(providers, accuProvider, openProvider, climacellProvider)

	response.Status = resolveStatus(providers)

	mintemp, _ := utils.Min(tempData)
	maxtemp, _ := utils.Max(tempData)

	desv, tmp, rem := utils.EnhancedTempCalculation(tempData)

	if response.Status != ERROR_STATUS {
		weather := models.Weather{
			MeanTemp:            tmp,
			StandardDeviation:   desv,
			WeatherProviders:    providers,
			MaxTemp:             maxtemp,
			MinTemp:             mintemp,
			SensorsDesestimated: rem,
			KnownTotalSensors:   len(tempData),
		}
		response.Payload = weather
	}

	return response
}

func resolveStatus(providers []models.WeatherProvider) string {
	okCount := 0
	for _, weatherProvider := range providers {
		if weatherProvider.Status == ONLINE_STATUS {
			okCount++
		}
	}
	if okCount == len(providers) {
		return OK_STATUS
	} else if okCount == 0 {
		return ERROR_STATUS
	} else {
		return WARNING_STATUS
	}
}

func getLocationTemperatureFromAccuweatherAPI(ctx context.Context, geo *client.Geoposition, provider *models.WeatherProvider, channel chan standardWeatherInfo, accuweatherProvider *AccuweatherClientWrapper) {
	defer waitGroup.Done()
	city, err := accuweatherProvider.GetAccuweatherCityByGeoposition(ctx, geo)
	if err != nil {
		setError(provider, err)
	}
	clima, err := accuweatherProvider.GetAccuweatherCurrentWeatherByCityKey(ctx, city.Key)
	if err != nil {
		setError(provider, err)
	}
	channel <- standardWeatherInfo{temperature: clima.Temperature.Metric.Value}
}

func getWeatherFromOpenWeatherAPI(ctx context.Context, geo *client.Geoposition, provider *models.WeatherProvider, channel chan standardWeatherInfo, openweathermapProvider *OpenweatherClientWrapper) {
	defer waitGroup.Done()
	clima, err := openweathermapProvider.GetOpenweatherMapWeather(ctx, geo)
	if err != nil {
		setError(provider, err)
	}
	channel <- standardWeatherInfo{temperature: clima.Main.Temp}
}

func getWeatherFromClimacellAPI(ctx context.Context, geo *client.Geoposition, provider *models.WeatherProvider, channel chan standardWeatherInfo, climacellProvider *ClimacellClientWrapper) {
	defer waitGroup.Done()
	clima, err := climacellProvider.GetClimacellWeatherCells(ctx, geo)
	if err != nil {
		setError(provider, err)
	}
	for _, cell := range *clima {
		channel <- standardWeatherInfo{temperature: cell.Temp.Value}
	}
}

func setError(provider *models.WeatherProvider, err error) {
	provider.Status = OFFLINE_STATUS
	provider.Error = err.Error()
}
