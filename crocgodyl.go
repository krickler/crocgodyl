package crocgodyl

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

// VERSION of crocgodyl follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.0.2-alpha"

// ErrorResponse is the response from the panels for errors.
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error is the struct from the panels for errors.
type Error struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
	Source struct {
		Field string `json:"field,omitempty"`
	} `json:"source,omitempty"`
}

// Meta is the meta information on some queries
type Meta struct {
	Pagination Pagination `json:"pagination,omitempty"`
}

// Pagination is the information how many responses there on a page and how many pages there are.
type Pagination struct {
	Total       int   `json:"total,omitempty"`
	Count       int   `json:"count,omitempty"`
	PerPage     int   `json:"per_page,omitempty"`
	CurrentPage int   `json:"current_page,omitempty"`
	TotalPages  int   `json:"total_pages,omitempty"`
	Links       Links `json:"links"`
}

type Links struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

//
// crocgodyl
//

// CrocConfig is the config for crocgodyl
type CrocConfig struct {
	PanelURL    string
	ClientToken string
	AppToken    string
}

//
// Application code
//

// NewCrocConfig sets up the API interface with
func NewApplication(panelURL string, appToken string) (config *CrocConfig, err error) {

	config = &CrocConfig{}

	if panelURL == "" && appToken == "" {
		return config, errors.New("you need to configure the panel and at least one api token")
	}

	if panelURL == "" {
		return config, errors.New("a panel URL is required to use the API")
	}

	if appToken == "" {
		return config, errors.New("an application token is required")
	}

	config.PanelURL = panelURL
	config.AppToken = appToken

	// validate the server is up and available
	if _, err = config.GetUsers(); err != nil {
		return config, err
	}

	return
}

func (config *CrocConfig) queryApplicationAPI(endpoint, request string, data []byte) (bodyBytes []byte, err error) {
	//http get json request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", config.PanelURL+"/api/application/"+endpoint, nil)

	switch {
	case request == "get":
	case request == "post":
		req, _ = http.NewRequest("POST", config.PanelURL+"/api/application/"+endpoint, bytes.NewBuffer(data))
	case request == "patch":
		req, _ = http.NewRequest("PATCH", config.PanelURL+"/api/application/"+endpoint, bytes.NewBuffer(data))
	case request == "delete":
		req, _ = http.NewRequest("DELETE", config.PanelURL+"/api/application/"+endpoint, nil)
	default:
	}

	//Sets request header for the http request
	req.Header.Add("Authorization", "Bearer "+config.AppToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 204 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New(string(bodyBytes))
	}

	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	//Close response thread
	defer resp.Body.Close()

	bodyBytes = bytes.Replace(bodyBytes, []byte("[]"), []byte("{}"), -1)

	//return byte structure
	return
}
