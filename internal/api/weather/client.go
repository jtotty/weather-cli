package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.weatherapi.com/v1/forecast.json"
const maxResponseSize = 10 * 1024 * 1024 // 10MB to prevent DoS

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:          10,
				IdleConnTimeout:       30 * time.Second,
				DisableCompression:    false,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
			},
		},
		apiKey: apiKey,
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

	body, err := io.ReadAll(io.LimitReader(res.Body, maxResponseSize))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if len(body) == maxResponseSize {
		return nil, fmt.Errorf("response too large (exceeded %d bytes)", maxResponseSize)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &response, nil
}

// Format URL preventing injection and properly encodes special chars
func (c *Client) buildURL(opts FetchOptions) string {
	params := url.Values{}
	params.Add("key", c.apiKey)
	params.Add("q", opts.Location)
	params.Add("days", fmt.Sprintf("%d", opts.Days))

	if opts.IncludeAQI {
		params.Add("aqi", "yes")
	}

	if opts.Alerts {
		params.Add("alerts", "yes")
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}
