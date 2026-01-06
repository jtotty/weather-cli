package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jtotty/weather-cli/internal/api/weather"
	"github.com/jtotty/weather-cli/internal/config"
)

// mockCache implements WeatherCache for testing.
type mockCache struct {
	data     map[string]*weather.Response
	getCalls []string
	setCalls []setCacheCall
	setError error
}

type setCacheCall struct {
	location string
	data     *weather.Response
}

func newMockCache() *mockCache {
	return &mockCache{
		data: make(map[string]*weather.Response),
	}
}

func (m *mockCache) Get(location string) *weather.Response {
	m.getCalls = append(m.getCalls, location)
	return m.data[location]
}

func (m *mockCache) Set(location string, data *weather.Response) error {
	m.setCalls = append(m.setCalls, setCacheCall{location, data})
	if m.setError != nil {
		return m.setError
	}
	m.data[location] = data
	return nil
}

// mockFetcher implements WeatherFetcher for testing.
type mockFetcher struct {
	response   *weather.Response
	err        error
	fetchCalls []weather.FetchOptions
}

func (m *mockFetcher) Fetch(ctx context.Context, opts weather.FetchOptions) (*weather.Response, error) {
	m.fetchCalls = append(m.fetchCalls, opts)
	return m.response, m.err
}

func TestNewWeather(t *testing.T) {
	cfg := &config.Config{
		APIKey:     "test-key",
		Location:   "London",
		Days:       3,
		IncludeAQI: true,
		Alerts:     true,
	}

	svc := NewWeather(cfg)

	if svc == nil {
		t.Fatal("NewWeather() returned nil")
	}

	if svc.cfg != cfg {
		t.Error("NewWeather() did not store config reference")
	}
}

func TestGetWeather_CacheHit(t *testing.T) {
	cfg := &config.Config{
		APIKey:   "test-key",
		Location: "London",
		Days:     3,
	}

	cachedResponse := &weather.Response{
		Location: weather.Location{Name: "London", Country: "UK"},
		Current:  weather.Current{TempC: 15},
	}

	mockCache := newMockCache()
	mockCache.data["London"] = cachedResponse

	mockFetcher := &mockFetcher{}

	svc := NewWeatherWithDeps(cfg, mockCache, mockFetcher)

	result, err := svc.GetWeather(context.Background())
	if err != nil {
		t.Fatalf("GetWeather() error = %v", err)
	}

	if result != cachedResponse {
		t.Error("GetWeather() did not return cached response")
	}

	if len(mockCache.getCalls) != 1 || mockCache.getCalls[0] != "London" {
		t.Errorf("Expected cache.Get(\"London\"), got %v", mockCache.getCalls)
	}

	if len(mockFetcher.fetchCalls) != 0 {
		t.Error("GetWeather() should not call API when cache hit")
	}
}

func TestGetWeather_CacheMiss(t *testing.T) {
	cfg := &config.Config{
		APIKey:     "test-key",
		Location:   "Paris",
		Days:       5,
		IncludeAQI: true,
		Alerts:     true,
	}

	apiResponse := &weather.Response{
		Location: weather.Location{Name: "Paris", Country: "France"},
		Current:  weather.Current{TempC: 20},
	}

	mockCache := newMockCache() // Empty cache
	mockFetcher := &mockFetcher{response: apiResponse}

	svc := NewWeatherWithDeps(cfg, mockCache, mockFetcher)

	result, err := svc.GetWeather(context.Background())
	if err != nil {
		t.Fatalf("GetWeather() error = %v", err)
	}

	if result != apiResponse {
		t.Error("GetWeather() did not return API response")
	}

	if len(mockCache.getCalls) != 1 {
		t.Errorf("Expected 1 cache.Get call, got %d", len(mockCache.getCalls))
	}

	if len(mockFetcher.fetchCalls) != 1 {
		t.Fatalf("Expected 1 API call, got %d", len(mockFetcher.fetchCalls))
	}

	opts := mockFetcher.fetchCalls[0]
	if opts.Location != "Paris" {
		t.Errorf("Fetch Location = %q, want %q", opts.Location, "Paris")
	}
	if opts.Days != 5 {
		t.Errorf("Fetch Days = %d, want %d", opts.Days, 5)
	}
	if !opts.IncludeAQI {
		t.Error("Fetch IncludeAQI should be true")
	}
	if !opts.Alerts {
		t.Error("Fetch Alerts should be true")
	}

	if len(mockCache.setCalls) != 1 {
		t.Fatalf("Expected 1 cache.Set call, got %d", len(mockCache.setCalls))
	}
	if mockCache.setCalls[0].location != "Paris" {
		t.Errorf("cache.Set location = %q, want %q", mockCache.setCalls[0].location, "Paris")
	}
	if mockCache.setCalls[0].data != apiResponse {
		t.Error("cache.Set did not receive correct data")
	}
}

func TestGetWeather_APIError(t *testing.T) {
	cfg := &config.Config{
		APIKey:   "test-key",
		Location: "InvalidLocation",
		Days:     1,
	}

	mockCache := newMockCache()
	mockFetcher := &mockFetcher{err: errors.New("API error: location not found")}

	svc := NewWeatherWithDeps(cfg, mockCache, mockFetcher)

	result, err := svc.GetWeather(context.Background())

	if err == nil {
		t.Fatal("GetWeather() expected error, got nil")
	}

	if result != nil {
		t.Error("GetWeather() expected nil result on error")
	}

	if len(mockCache.setCalls) != 0 {
		t.Error("GetWeather() should not cache on API error")
	}
}

func TestGetWeather_CacheSetError(t *testing.T) {
	cfg := &config.Config{
		APIKey:   "test-key",
		Location: "London",
		Days:     1,
	}

	apiResponse := &weather.Response{
		Location: weather.Location{Name: "London"},
	}

	mockCache := newMockCache()
	mockCache.setError = errors.New("cache write failed")

	mockFetcher := &mockFetcher{response: apiResponse}

	svc := NewWeatherWithDeps(cfg, mockCache, mockFetcher)

	result, err := svc.GetWeather(context.Background())
	if err != nil {
		t.Fatalf("GetWeather() error = %v", err)
	}

	if result != apiResponse {
		t.Error("GetWeather() should return data even when cache fails")
	}
}

func TestGetWeather_NilCache(t *testing.T) {
	cfg := &config.Config{
		APIKey:   "test-key",
		Location: "Tokyo",
		Days:     1,
	}

	apiResponse := &weather.Response{
		Location: weather.Location{Name: "Tokyo"},
	}

	mockFetcher := &mockFetcher{response: apiResponse}

	svc := NewWeatherWithDeps(cfg, nil, mockFetcher)

	result, err := svc.GetWeather(context.Background())
	if err != nil {
		t.Fatalf("GetWeather() error = %v", err)
	}

	if result != apiResponse {
		t.Error("GetWeather() should work without cache")
	}

	if len(mockFetcher.fetchCalls) != 1 {
		t.Error("GetWeather() should call API when no cache")
	}
}

func TestGetWeather_ContextCancellation(t *testing.T) {
	cfg := &config.Config{
		APIKey:   "test-key",
		Location: "Berlin",
		Days:     1,
	}

	mockCache := newMockCache()
	mockFetcher := &mockFetcher{err: context.Canceled}

	svc := NewWeatherWithDeps(cfg, mockCache, mockFetcher)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.GetWeather(ctx)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("GetWeather() error = %v, want context.Canceled", err)
	}
}
