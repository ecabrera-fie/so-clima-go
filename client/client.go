package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	LAT_DEFAULT = "-34.6083"
	LON_DEFAULT = "-58.3712"
)

type Geoposition struct {
	Latitude  string
	Longitude string
}

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

type AccuweatherCurrentWeather []struct {
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

type Temperature struct {
	Metric Metric `json:"Metric"`
}

type Metric struct {
	Value string `json:"Value"`
}

type Client struct {
	Name    string
	BaseURL string
	apiKey  string
}

func NewClient(name string, apiKey string, baseUrl string) *Client {
	return &Client{
		Name:    name,
		BaseURL: baseUrl,
		apiKey:  apiKey,
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetAccuweatherCityByGeoposition obtiene una ciudad a través de geolocalización
func (c *Client) GetAccuweatherCityByGeoposition(ctx context.Context, geoposition *Geoposition) (*AccuweatherCity, error) {
	lat := LAT_DEFAULT
	lon := LON_DEFAULT
	if *geoposition != (Geoposition{}) {
		lat = geoposition.Latitude
		lon = geoposition.Longitude
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/locations/v1/cities/geoposition/search?apikey=%s&q=%s,%s", c.BaseURL, c.apiKey, lat, lon), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := AccuweatherCity{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAccuweatherCurrentWeatherByCityKey A partir de una clave de ciudad, obtiene el clima
func (c *Client) GetAccuweatherCurrentWeatherByCityKey(ctx context.Context, key string) (*AccuweatherCurrentWeather, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/currentconditions/v1/%s?apikey=%s", c.BaseURL, key, c.apiKey), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := AccuweatherCurrentWeather{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	respBody, err := ioutil.ReadAll(res.Body)

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.Unmarshal(respBody, &v); err != nil {
		return err
	}

	return nil
}
