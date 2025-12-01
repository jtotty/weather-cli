package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jtotty/weather-cli/internal/api/weather"
)

const (
	DefaultTTL = 30 * time.Minute
	cacheDir   = "weather-cli"
	cacheFile  = "cache.json"
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
	path    string
	ttl     time.Duration
}

func New(ttl time.Duration) (*Cache, error) {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	cachePath := filepath.Join(cacheDir, cacheFile)

	cache := &Cache{
		Entries: make(map[string]*Entry),
		path:    cachePath,
		ttl:     ttl,
	}

	if err := cache.load(); err != nil && !os.IsNotExist(err) {
		return cache, nil
	}

	return cache, nil
}

// Get retrieves a cached weather response for the given location.
func (c *Cache) Get(location string) *weather.Response {
	key := normalizeKey(location)
	entry, ok := c.Entries[key]
	if !ok {
		return nil
	}

	if !entry.IsValid(c.ttl) {
		delete(c.Entries, key)
		return nil
	}

	return entry.Data
}

// Set stores a weather response in the cache for the given location.
func (c *Cache) Set(location string, data *weather.Response) error {
	key := normalizeKey(location)
	c.Entries[key] = &Entry{
		Location: location,
		Data:     data,
		CachedAt: time.Now(),
	}

	return c.save()
}

func (c *Cache) Clear() error {
	c.Entries = make(map[string]*Entry)
	return c.save()
}

func (c *Cache) Path() string {
	return c.path
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
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(c.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

func getCacheDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userCacheDir, cacheDir), nil
}

// normalizeKey creates a consistent cache key from a location string.
// Uses SHA256 hash to handle special characters and long location names.
func normalizeKey(location string) string {
	normalized := strings.ToLower(strings.TrimSpace(location))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:8]) // Use first 8 bytes (16 hex chars)
}
