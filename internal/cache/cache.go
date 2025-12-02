package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/jtotty/weather-cli/internal/api/weather"
)

const (
	DefaultTTL      = 30 * time.Minute
	cacheSubDir     = "weather-cli"
	cacheFileName   = "cache.json"
	maxCacheEntries = 100
)

type Entry struct {
	Location string            `json:"location"`
	Data     *weather.Response `json:"data"`
	CachedAt time.Time         `json:"cached_at"`
}

func (e *Entry) IsValid(ttl time.Duration) bool {
	return time.Since(e.CachedAt) < ttl
}

type Cache struct {
	Entries map[string]*Entry `json:"entries"`
	path    string            `json:"-"`
	ttl     time.Duration     `json:"-"`
	mu      sync.RWMutex      `json:"-"`
}

func New(ttl time.Duration) (*Cache, error) {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	cachePath := filepath.Join(cacheDir, cacheFileName)

	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    cachePath,
		ttl:     ttl,
	}

	if err := cache.load(); err != nil {
		if os.IsNotExist(err) {
			return cache, nil
		}
		return cache, nil
	}

	return cache, nil
}

func (c *Cache) Get(location string) *weather.Response {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := normalizeKey(location)
	entry, ok := c.Entries[key]
	if !ok {
		return nil
	}

	if !entry.IsValid(c.ttl) {
		return nil
	}

	return entry.Data
}

// getCacheDir returns the cache directory for weather-cli.
func (c *Cache) Set(location string, data *weather.Response) error {
	if data == nil {
		return errors.New("cannot cache nil weather data")
	}

	location = strings.TrimSpace(location)
	if location == "" {
		return errors.New("cannot cache empty location")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.cleanupExpired()

	if len(c.Entries) >= maxCacheEntries {
		c.removeOldest()
	}

	key := normalizeKey(location)
	c.Entries[key] = &Entry{
		Location: location,
		Data:     data,
		CachedAt: time.Now().UTC(),
	}

	return c.save()
}

func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries = make(map[string]*Entry)
	return c.save()
}

func (c *Cache) Path() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.path
}

func (c *Cache) Stats() (total, valid, expired int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total = len(c.Entries)
	for _, entry := range c.Entries {
		if entry.IsValid(c.ttl) {
			valid++
		} else {
			expired++
		}
	}
	return
}

func (c *Cache) load() error {
	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

func (c *Cache) save() error {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	// Atomic write: write to temp file first, then rename
	tmpFile := c.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o600); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, c.path); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("failed to rename cache file: %w", err)
	}

	return nil
}

func (c *Cache) cleanupExpired() {
	for key, entry := range c.Entries {
		if !entry.IsValid(c.ttl) {
			delete(c.Entries, key)
		}
	}
}

func (c *Cache) removeOldest() {
	var oldestKey string
	var oldestTime time.Time

	first := true
	for key, entry := range c.Entries {
		if first || entry.CachedAt.Before(oldestTime) {
			oldestTime = entry.CachedAt
			oldestKey = key
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.Entries, oldestKey)
	}
}

func getCacheDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userCacheDir, cacheSubDir), nil
}

// normalizeKey creates a consistent cache key from a location string.
// Uses lowercase and trimmed string for case-insensitive matching.
func normalizeKey(location string) string {
	return strings.ToLower(strings.TrimSpace(location))
}
