package providers

import (
	"context"
	"so2-clima/client"
	"so2-clima/models"
	"so2-clima/utils"
	"sync"
)

var waitGroup sync.WaitGroup

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
	accuChannel	chan *AccuweatherCurrentWeather
	openweathermapProvider *OpenweatherClientWrapper
	openChannel	chan *OpenweatherMapWeather
	climacellProvider *ClimacellClientWrapper
	climacellChannel chan *[]ClimacellCell
}

func NewDistributedWeatherProvider(ctx context.Context) *ProvidersFacade {
	return &ProvidersFacade{
		Context:                ctx,
		accuweatherProvider:    NewAccuweatherClient(),
		accuChannel: make(chan *AccuweatherCurrentWeather),
		openweathermapProvider: NewOpenweatherMapClient(),
		openChannel: make(chan *OpenweatherMapWeather),
		climacellProvider: NewClimacellClient(),
		climacellChannel: make(chan *[]ClimacellCell),
	}
}

func (p *ProvidersFacade) GetTemperatureDataByGeolocation(geo *client.Geoposition) models.Response {
	var sumTemp float64
	var tempData []float64
	var providers []models.WeatherProvider

	response := models.Response{}

	waitGroup.Add(3)

	go func() {
		waitGroup.Wait()
		close(p.accuChannel)
		close(p.openChannel)
		close(p.climacellChannel)
	}()

	accuProvider := models.WeatherProvider{ Name: p.accuweatherProvider.C.Name, Status: ONLINE_STATUS }
	go getLocationTemperatureFromAccuweatherAPI(p.Context, geo, &accuProvider, p.accuChannel, p.accuweatherProvider)

	openProvider := models.WeatherProvider{ Name: p.openweathermapProvider.C.Name, Status: ONLINE_STATUS }
	go getWeatherFromOpenWeatherAPI(p.Context, geo, &openProvider, p.openChannel, p.openweathermapProvider)

	climacellProvider := models.WeatherProvider{ Name: p.climacellProvider.C.Name, Status: ONLINE_STATUS }
	go getWeatherFromClimacellAPI(p.Context, geo, &climacellProvider, p.climacellChannel, p.climacellProvider)

	accu := <- p.accuChannel
	if accu != nil {
		tempData = append(tempData, accu.Temperature.Metric.Value)
		sumTemp += accu.Temperature.Metric.Value
	}

	open := <- p.openChannel
	if open != nil {
		tempData = append(tempData, open.Main.Temp)
		sumTemp += open.Main.Temp
	}

	cells := <- p.climacellChannel
	if cells != nil {
		for _, cell := range *cells {
			tempData = append(tempData, cell.Temp.Value)
			sumTemp += cell.Temp.Value
		}
	}

	providers = append(providers, accuProvider, openProvider, climacellProvider)

	okCount := 0
	for _, weatherProvider := range providers {
		if weatherProvider.Status == ONLINE_STATUS {
			okCount++
		}
	}
	if okCount == len(providers) {
		response.Status = OK_STATUS
	} else if okCount == 0 {
		response.Status = ERROR_STATUS
	} else {
		response.Status = WARNING_STATUS
	}

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

func getLocationTemperatureFromAccuweatherAPI(ctx context.Context, geo *client.Geoposition, provider *models.WeatherProvider, channel chan *AccuweatherCurrentWeather, accuweatherProvider *AccuweatherClientWrapper) {
	defer waitGroup.Done()
	city, err := accuweatherProvider.GetAccuweatherCityByGeoposition(ctx, geo)
	if err != nil {
		setError(provider, err)
	}
	clima, err := accuweatherProvider.GetAccuweatherCurrentWeatherByCityKey(ctx, city.Key)
	if err != nil {
		setError(provider, err)
	}
	channel <- clima
}

func getWeatherFromOpenWeatherAPI(ctx context.Context, geo *client.Geoposition, provider *models.WeatherProvider, channel chan *OpenweatherMapWeather, openweathermapProvider *OpenweatherClientWrapper) {
	defer waitGroup.Done()
	clima, err := openweathermapProvider.GetOpenweatherMapWeather(ctx, geo)
	if err != nil {
		setError(provider, err)
	}
	channel <- clima
}

func getWeatherFromClimacellAPI(ctx context.Context, geo *client.Geoposition, provider *models.WeatherProvider, channel chan *[]ClimacellCell, climacellProvider *ClimacellClientWrapper) {
	defer waitGroup.Done()
	clima, err := climacellProvider.GetClimacellWeatherCells(ctx, geo)
	if err != nil {
		setError(provider, err)
	}
	channel <- clima
}

func setError(provider *models.WeatherProvider, err error)  {
	provider.Status = OFFLINE_STATUS
	provider.Error = err.Error()
}