package providers

import (
	"context"
	"fmt"
	"net/http"
	"so2-clima/client"
	"so2-clima/constants"
)

type AccuweatherCity struct {
	Version           int    `json:"Version"`
	Key               string `json:"Key"`
	Type              string `json:"Type"`
	Rank              int    `json:"Rank"`
	LocalizedName     string `json:"LocalizedName"`
	EnglishName       string `json:"EnglishName"`
	PrimaryPostalCode string `json:"PrimaryPostalCode"`
	Region            struct {
		ID            string `json:"ID"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"Region"`
	Country struct {
		ID            string `json:"ID"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"Country"`
	AdministrativeArea struct {
		ID            string `json:"ID"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
		Level         int    `json:"Level"`
		LocalizedType string `json:"LocalizedType"`
		EnglishType   string `json:"EnglishType"`
		CountryID     string `json:"CountryID"`
	} `json:"AdministrativeArea"`
	TimeZone struct {
		Code             string      `json:"Code"`
		Name             string      `json:"Name"`
		GmtOffset        float64     `json:"GmtOffset"`
		IsDaylightSaving bool        `json:"IsDaylightSaving"`
		NextOffsetChange interface{} `json:"NextOffsetChange"`
	} `json:"TimeZone"`
	GeoPosition struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
		Elevation struct {
			Metric struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Metric"`
			Imperial struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Imperial"`
		} `json:"Elevation"`
	} `json:"GeoPosition"`
	IsAlias    bool `json:"IsAlias"`
	ParentCity struct {
		Key           string `json:"Key"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"ParentCity"`
	SupplementalAdminAreas []struct {
		Level         int    `json:"Level"`
		LocalizedName string `json:"LocalizedName"`
		EnglishName   string `json:"EnglishName"`
	} `json:"SupplementalAdminAreas"`
	DataSets []string `json:"DataSets"`
}

type AccuweatherCurrentWeather struct {
	LocalObservationDateTime string      `json:"LocalObservationDateTime"`
	EpochTime                int         `json:"EpochTime"`
	WeatherText              string      `json:"WeatherText"`
	WeatherIcon              int         `json:"WeatherIcon"`
	HasPrecipitation         bool        `json:"HasPrecipitation"`
	PrecipitationType        interface{} `json:"PrecipitationType"`
	IsDayTime                bool        `json:"IsDayTime"`
	Temperature              struct {
		Metric struct {
			Value    float64 `json:"Value"`
			Unit     string  `json:"Unit"`
			UnitType int     `json:"UnitType"`
		} `json:"Metric"`
		Imperial struct {
			Value    float64 `json:"Value"`
			Unit     string  `json:"Unit"`
			UnitType int     `json:"UnitType"`
		} `json:"Imperial"`
	} `json:"Temperature"`
	MobileLink string `json:"MobileLink"`
	Link       string `json:"Link"`
}

type AccuweatherClientWrapper struct {
	C *client.Client
}

func NewAccuweatherClient() *AccuweatherClientWrapper {
	return &AccuweatherClientWrapper{
		C: client.NewClient(constants.ACCUWEATHER_API_NAME, constants.ACCUWEATHER_API_KEY, constants.ACCUWEATHER_BASE_PATH),
	}
}

// GetAccuweatherCityByGeoposition Obtiene una ciudad a través de geolocalización
func (acw *AccuweatherClientWrapper) GetAccuweatherCityByGeoposition(ctx context.Context, geoposition *client.Geoposition) (*AccuweatherCity, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/locations/v1/cities/geoposition/search?apikey=%s&q=%s,%s", acw.C.BaseURL, acw.C.ApiKey, geoposition.Latitude, geoposition.Longitude), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := AccuweatherCity{}

	if err := acw.C.SendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAccuweatherCurrentWeatherByCityKey A partir de una clave de ciudad, obtiene el clima
func (acw *AccuweatherClientWrapper) GetAccuweatherCurrentWeatherByCityKey(ctx context.Context, key string) (*AccuweatherCurrentWeather, error) {
	var res []AccuweatherCurrentWeather
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/currentconditions/v1/%s?apikey=%s", acw.C.BaseURL, key, acw.C.ApiKey), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	if err := acw.C.SendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res[0], nil
}
