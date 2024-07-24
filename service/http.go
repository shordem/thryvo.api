package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

var client = &http.Client{}

type HttpService struct{}

type HttpServiceInterface interface {
	Get(url string, headers map[string]string) (*http.Response, error)
	Post(url string, headers map[string]string, body interface{}) (*http.Response, error)
	BodyToDTO(body io.ReadCloser, dto interface{}) error
}

func NewHTTPService() HttpServiceInterface {
	return &HttpService{}
}

func (service *HttpService) Get(url string, headers map[string]string) (*http.Response, error) {
	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *HttpService) Post(url string, headers map[string]string, body interface{}) (*http.Response, error) {
	// Create the request body as JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *HttpService) BodyToDTO(body io.ReadCloser, dto interface{}) error {
	return json.NewDecoder(body).Decode(dto)
}
