package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://api.weatherapi.com/v1/forecast.json"

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: http.DefaultClient,
		apiKey:     apiKey,
	}
}

type FetchOptions struct {
	Location   string
	Days       int
	IncludeAQI bool
	Alerts     bool
}

func (c *Client) Fetch(opts FetchOptions) (*Response, error) {
	url := c.buildURL(opts)

	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &response, nil
}

func (c *Client) buildURL(opts FetchOptions) string {
	url := fmt.Sprintf("%s?key=%s&q=%s&days=%d", baseURL, c.apiKey, opts.Location, opts.Days)

	if opts.IncludeAQI {
		url += "&aqi=yes"
	}

	if opts.Alerts {
		url += "&alerts=yes"
	}

	return url
}
