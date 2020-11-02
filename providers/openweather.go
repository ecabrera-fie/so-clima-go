package providers

import (
	"context"
	"fmt"
	"net/http"
	"so2-clima/client"
	"so2-clima/constants"
)

type OpenweatherMapWeather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type OpenweatherClientWrapper struct {
	C *client.Client
}

func NewOpenweatherMapClient() *OpenweatherClientWrapper {
	return &OpenweatherClientWrapper{
		C: client.NewClient(constants.OPENWEATHER_API_NAME, constants.OPENWEATHER_API_KEY, constants.OPENWEATHER_BASE_PATH),
	}
}

func (owc *OpenweatherClientWrapper) GetOpenweatherMapWeather(ctx context.Context, geoposition *client.Geoposition) (*OpenweatherMapWeather, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/weather?lat=%s&lon=%s&appid=%s&units=metric", owc.C.BaseURL, geoposition.Latitude, geoposition.Longitude, owc.C.ApiKey), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := OpenweatherMapWeather{}

	if err := owc.C.SendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
