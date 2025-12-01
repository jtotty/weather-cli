package weather

import (
	"context"
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
	baseURL    string
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
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

type FetchOptions struct {
	Location   string
	Days       int
	IncludeAQI bool
	Alerts     bool
}

func (c *Client) Fetch(ctx context.Context, opts FetchOptions) (*Response, error) {
	reqURL := c.buildURL(opts)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

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

// buildURL constructs the API URL with proper encoding to prevent injection
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

	return fmt.Sprintf("%s?%s", c.baseURL, params.Encode())
}
