package toggl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
)

const baseTogglUrl = "https://api.track.toggl.com/"
const basicAuthPassword = "api_token"
const mePath = "api/v9/me"

type TogglClient struct {
	baseUrl    url.URL
	httpClient *http.Client
	apiToken   string
}

func NewTogglClient(apiToken string) TogglClient {
	baseUrl, _ := url.Parse(baseTogglUrl)
	client := TogglClient{
		baseUrl:    *baseUrl,
		httpClient: http.DefaultClient,
		apiToken:   apiToken,
	}

	return client
}

func (c TogglClient) httpGet(urlPath string) (*http.Response, error) {
	c.baseUrl.Path = path.Join(c.baseUrl.Path, urlPath)
	req, _ := http.NewRequest(http.MethodGet, c.baseUrl.String(), nil)
	req.SetBasicAuth(c.apiToken, basicAuthPassword)
	req.Header.Set("content-type", "application/json")

	return c.httpClient.Do(req)
}

func (c TogglClient) httpPost(urlPath string, body any) (*http.Response, error) {
	c.baseUrl.Path = path.Join(c.baseUrl.Path, urlPath)
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	if err := encoder.Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseUrl.String(), &buffer)

	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.apiToken, basicAuthPassword)
	return c.httpClient.Do(req)

}

func (c TogglClient) httpPut(urlPath string, body any) (*http.Response, error) {
	c.baseUrl.Path = path.Join(c.baseUrl.Path, urlPath)
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	if err := encoder.Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, c.baseUrl.String(), &buffer)

	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.apiToken, basicAuthPassword)
	return c.httpClient.Do(req)

}
