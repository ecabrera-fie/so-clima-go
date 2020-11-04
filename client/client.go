package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Geoposition struct {
	Latitude  string
	Longitude string
}

type Client struct {
	Name    string
	BaseURL string
	ApiKey  string
}

func NewClient(name string, apiKey string, baseUrl string) *Client {
	return &Client{
		Name:    name,
		BaseURL: baseUrl,
		ApiKey:  apiKey,
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) SendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	respBody, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(respBody, &v); err != nil {
		return err
	}

	return nil
}
