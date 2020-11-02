package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"so2-clima/client"
	"so2-clima/constants"
	"so2-clima/providers"
)

var c *providers.AccuweatherClientWrapper

// HandleGetTemperature handles getting the temperature
func HandleGetTemperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	lat := r.URL.Query()["lat"]
	lon := r.URL.Query()["lon"]

	geo := client.Geoposition{}

	if lat != nil {
		geo.Latitude = lat[0]
	} else {
		geo.Latitude = constants.DEFAULT_LATITUDE
	}

	if lon != nil {
		geo.Longitude = lon[0]
	} else {
		geo.Longitude = constants.DEFAULT_LONGITUDE
	}

	c = providers.NewAccuweatherClient()

	clima, err := getLocationTemperatureFromAccuweatherAPI(ctx, &geo)
	if err != nil {
		json.NewEncoder(w).Encode("")
	}

	json.NewEncoder(w).Encode(clima)
}

func getLocationTemperatureFromAccuweatherAPI(ctx context.Context, geo *client.Geoposition) (*providers.AccuweatherCurrentWeather, error) {
	city, err := c.GetAccuweatherCityByGeoposition(ctx, geo)
	if err != nil {
		return nil, err
	}

	clima, err := c.GetAccuweatherCurrentWeatherByCityKey(ctx, city.Key)
	if err != nil {
		return nil, err
	}

	return clima, nil
}
