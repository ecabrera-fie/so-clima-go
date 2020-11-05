package providers

import (
	"context"
	"fmt"
	"net/http"
	"so2-clima/client"
	"so2-clima/constants"
)

type ClimacellCell struct {
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Temp struct {
		Value float64 `json:"value"`
		Units string  `json:"units"`
	} `json:"temp"`
	ObservationTime struct {
		Value string `json:"value"`
	} `json:"observation_time"`
}

type ClimacellClientWrapper struct {
	C *client.Client
}

func NewClimacellClient() *ClimacellClientWrapper {
	return &ClimacellClientWrapper{
		C: client.NewClient(constants.CLIMACELL_API_NAME, constants.CLIMACELL_API_KEY, constants.CLIMACELL_BASE_PATH),
	}
}

// GetAccuweatherCityByGeoposition Obtiene una ciudad a través de geolocalización
func (ccw *ClimacellClientWrapper) GetClimacellWeatherCells(ctx context.Context, geoposition *client.Geoposition) (*[]ClimacellCell, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/weather/nowcast?apikey=%s&lat=%s&lon=%s&unit_system=si&timestep=5&start_time=now&fields=temp", ccw.C.BaseURL, ccw.C.ApiKey, geoposition.Latitude, geoposition.Longitude), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	var res []ClimacellCell

	if err := ccw.C.SendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
