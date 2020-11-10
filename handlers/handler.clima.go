package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"../client"
	"../constants"
	"../providers"
)

// HandleGetTemperature handles getting the temperature
func HandleGetTemperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := context.Background()

	lat := r.URL.Query()["lat"]
	lon := r.URL.Query()["lon"]

	geo := getGeoposition(lat, lon)

	provider := providers.NewDistributedWeatherProvider(ctx)

	clima := provider.GetTemperatureDataByGeolocation(geo)

	json.NewEncoder(w).Encode(clima)
}

func getGeoposition(lat []string, lon []string) *client.Geoposition {
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

	return &geo
}
