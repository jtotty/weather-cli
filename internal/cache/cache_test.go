package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jtotty/weather-cli/internal/api/weather"
)

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		name     string
		input1   string
		input2   string
		wantSame bool
	}{
		{
			name:     "same location different case",
			input1:   "London",
			input2:   "london",
			wantSame: true,
		},
		{
			name:     "same location with whitespace",
			input1:   "  London  ",
			input2:   "London",
			wantSame: true,
		},
		{
			name:     "different locations",
			input1:   "London",
			input2:   "Paris",
			wantSame: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key1 := normalizeKey(tt.input1)
			key2 := normalizeKey(tt.input2)

			if tt.wantSame && key1 != key2 {
				t.Errorf("expected same keys for %q and %q, got %q and %q", tt.input1, tt.input2, key1, key2)
			}
			if !tt.wantSame && key1 == key2 {
				t.Errorf("expected different keys for %q and %q, got same key %q", tt.input1, tt.input2, key1)
			}
		})
	}
}

func TestEntryIsValid(t *testing.T) {
	tests := []struct {
		name     string
		cachedAt time.Time
		ttl      time.Duration
		want     bool
	}{
		{
			name:     "fresh entry",
			cachedAt: time.Now(),
			ttl:      30 * time.Minute,
			want:     true,
		},
		{
			name:     "expired entry",
			cachedAt: time.Now().Add(-31 * time.Minute),
			ttl:      30 * time.Minute,
			want:     false,
		},
		{
			name:     "just before expiry",
			cachedAt: time.Now().Add(-29 * time.Minute),
			ttl:      30 * time.Minute,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := &Entry{CachedAt: tt.cachedAt}
			if got := entry.IsValid(tt.ttl); got != tt.want {
				t.Errorf("Entry.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCache(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create cache with short TTL for testing
	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    filepath.Join(tmpDir, "cache.json"),
		ttl:     1 * time.Second,
	}

	// Create mock weather response
	mockResponse := &weather.Response{
		Location: weather.Location{
			Name:    "London",
			Country: "UK",
		},
	}

	t.Run("set and get", func(t *testing.T) {
		err := cache.Set("London", mockResponse)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		got := cache.Get("London")
		if got == nil {
			t.Fatal("Get() returned nil, expected cached data")
		}

		if got.Location.Name != mockResponse.Location.Name {
			t.Errorf("Get() location = %v, want %v", got.Location.Name, mockResponse.Location.Name)
		}
	})

	t.Run("get with different case", func(t *testing.T) {
		got := cache.Get("LONDON")
		if got == nil {
			t.Fatal("Get() returned nil, expected cached data (case insensitive)")
		}
	})

	t.Run("get non-existent", func(t *testing.T) {
		got := cache.Get("Paris")
		if got != nil {
			t.Errorf("Get() = %v, want nil for non-existent location", got)
		}
	})

	t.Run("expired entry", func(t *testing.T) {
		// Wait for cache to expire
		time.Sleep(2 * time.Second)

		got := cache.Get("London")
		if got != nil {
			t.Errorf("Get() = %v, want nil for expired entry", got)
		}
	})

	t.Run("clear", func(t *testing.T) {
		// Add a fresh entry
		err := cache.Set("Paris", mockResponse)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}

		err = cache.Clear()
		if err != nil {
			t.Fatalf("Clear() error = %v", err)
		}

		got := cache.Get("Paris")
		if got != nil {
			t.Errorf("Get() after Clear() = %v, want nil", got)
		}
	})
}

func TestCachePersistence(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-persist-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cachePath := filepath.Join(tmpDir, "cache.json")

	// Create mock weather response
	mockResponse := &weather.Response{
		Location: weather.Location{
			Name:    "London",
			Country: "UK",
		},
	}

	// Create first cache instance and save data
	cache1 := &Cache{
		Entries: make(map[string]*Entry),
		path:    cachePath,
		ttl:     30 * time.Minute,
	}

	err = cache1.Set("London", mockResponse)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Create second cache instance and verify data persisted
	cache2 := &Cache{
		Entries: make(map[string]*Entry),
		path:    cachePath,
		ttl:     30 * time.Minute,
	}

	err = cache2.load()
	if err != nil {
		t.Fatalf("load() error = %v", err)
	}

	got := cache2.Get("London")
	if got == nil {
		t.Fatal("Get() returned nil, expected persisted data")
	}

	if got.Location.Name != mockResponse.Location.Name {
		t.Errorf("Get() location = %v, want %v", got.Location.Name, mockResponse.Location.Name)
	}
}
