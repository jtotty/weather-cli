package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
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

func TestCacheInputValidation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-validation-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    filepath.Join(tmpDir, "cache.json"),
		ttl:     30 * time.Minute,
	}

	t.Run("nil data", func(t *testing.T) {
		err := cache.Set("London", nil)
		if err == nil {
			t.Error("Set() expected error for nil data, got nil")
		}
	})

	t.Run("empty location", func(t *testing.T) {
		mockResponse := &weather.Response{
			Location: weather.Location{Name: "London"},
		}
		err := cache.Set("", mockResponse)
		if err == nil {
			t.Error("Set() expected error for empty location, got nil")
		}
	})

	t.Run("whitespace only location", func(t *testing.T) {
		mockResponse := &weather.Response{
			Location: weather.Location{Name: "London"},
		}
		err := cache.Set("   ", mockResponse)
		if err == nil {
			t.Error("Set() expected error for whitespace-only location, got nil")
		}
	})
}

func TestCacheCorruptedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-corrupt-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cachePath := filepath.Join(tmpDir, "cache.json")

	// Write corrupted JSON
	err = os.WriteFile(cachePath, []byte("invalid json{{{"), 0600)
	if err != nil {
		t.Fatalf("failed to write corrupted cache: %v", err)
	}

	// Cache should still be usable with empty entries
	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    cachePath,
		ttl:     30 * time.Minute,
	}

	err = cache.load()
	if err == nil {
		t.Error("expected error when loading corrupted cache")
	}

	// Cache should still be usable
	if cache.Entries == nil {
		t.Error("cache.Entries should not be nil after failed load")
	}

	// Should be able to write new data
	mockResponse := &weather.Response{
		Location: weather.Location{Name: "London"},
	}
	err = cache.Set("London", mockResponse)
	if err != nil {
		t.Errorf("Set() after corrupted load failed: %v", err)
	}
}

func TestCacheStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-stats-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    filepath.Join(tmpDir, "cache.json"),
		ttl:     1 * time.Hour,
	}

	mockResponse := &weather.Response{
		Location: weather.Location{Name: "London"},
	}

	// Add some entries
	_ = cache.Set("London", mockResponse)
	_ = cache.Set("Paris", mockResponse)

	// Add an expired entry manually
	cache.mu.Lock()
	cache.Entries["expired"] = &Entry{
		Location: "Expired",
		Data:     mockResponse,
		CachedAt: time.Now().Add(-2 * time.Hour),
	}
	cache.mu.Unlock()

	total, valid, expired := cache.Stats()

	if total != 3 {
		t.Errorf("Stats() total = %d, want 3", total)
	}
	if valid != 2 {
		t.Errorf("Stats() valid = %d, want 2", valid)
	}
	if expired != 1 {
		t.Errorf("Stats() expired = %d, want 1", expired)
	}
}

func TestCacheMaxEntries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-max-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    filepath.Join(tmpDir, "cache.json"),
		ttl:     1 * time.Hour,
	}

	mockResponse := &weather.Response{
		Location: weather.Location{Name: "Test"},
	}

	// Add more than maxCacheEntries
	for i := 0; i < maxCacheEntries+10; i++ {
		location := fmt.Sprintf("Location%d", i)
		err := cache.Set(location, mockResponse)
		if err != nil {
			t.Fatalf("Set() error = %v", err)
		}
	}

	// Cache should not exceed maxCacheEntries
	total, _, _ := cache.Stats()
	if total > maxCacheEntries {
		t.Errorf("Cache size = %d, want <= %d", total, maxCacheEntries)
	}
}

func TestCacheConcurrency(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-race-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    filepath.Join(tmpDir, "cache.json"),
		ttl:     30 * time.Minute,
	}

	mockResponse := &weather.Response{
		Location: weather.Location{
			Name:    "London",
			Country: "UK",
		},
	}

	// Run with: go test -race
	const goroutines = 10
	const iterations = 50

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Writers
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				location := fmt.Sprintf("Location%d", id)
				_ = cache.Set(location, mockResponse)
			}
		}(i)
	}

	// Readers
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				location := fmt.Sprintf("Location%d", id)
				_ = cache.Get(location)
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is still usable after concurrent access
	got := cache.Get("Location0")
	if got == nil {
		t.Error("cache should have Location0 after concurrent writes")
	}
}

func TestCacheAtomicWrite(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "weather-cli-cache-atomic-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	cachePath := filepath.Join(tmpDir, "cache.json")

	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    cachePath,
		ttl:     30 * time.Minute,
	}

	mockResponse := &weather.Response{
		Location: weather.Location{Name: "London"},
	}

	err = cache.Set("London", mockResponse)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Verify no temp file remains
	tmpFile := cachePath + ".tmp"
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("temp file should not exist after successful write")
	}

	// Verify cache file exists with correct permissions
	info, err := os.Stat(cachePath)
	if err != nil {
		t.Fatalf("cache file should exist: %v", err)
	}

	// Check permissions (0600 = -rw-------)
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("cache file permissions = %o, want 0600", perm)
	}
}
