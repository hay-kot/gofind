package gofind

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var ErrCacheNotFound = errors.New("cache not found")

type Cache struct {
	Dir string
}

func (ce Cache) PathConstructor(namespace string) (string, error) {
	root, err := ParsePath(ce.Dir)
	if err != nil {
		return "", err
	}

	return filepath.Join(root, fmt.Sprintf("%s.json", namespace)), nil
}

func NewCache(dir string) (Cache, error) {
	p, err := ParsePath(dir)
	if err != nil {
		return Cache{}, err
	}

	// Create directory if it doesn't exist
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0755); err != nil {
			return Cache{}, fmt.Errorf("failed to create cache directory %s: %w", p, err)
		}
	}

	return Cache{
		Dir: dir,
	}, nil
}

// Find will find the first match in the cache or return an error if the cache
// isn't found.
func (c Cache) Find(namespace string) (*CacheEntry, error) {
	cachePath, err := c.PathConstructor(namespace)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, ErrCacheNotFound
	}

	// Load cache
	file, err := os.Open(cachePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close() // Ignore close errors for read-only operations
	}()

	var entry CacheEntry
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&entry); err != nil {
		return nil, err
	}

	// Check if cache is expired
	if entry.IsExpired() {
		return nil, ErrCacheNotFound
	}

	return &entry, nil
}

func (c Cache) Set(namespace string, results []Match) (*CacheEntry, error) {
	p, err := c.PathConstructor(namespace)
	if err != nil {
		return nil, err
	}

	entry := CacheEntry{
		Matches: results,
		Expires: time.Now().Add(time.Hour * 24),
	}

	// Create Parent
	err = os.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(p)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close() // Ignore close errors for read-only operations
	}()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(entry)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

type CacheEntry struct {
	Matches []Match   `json:"matches"`
	Expires time.Time `json:"expires"`
}

func (ce CacheEntry) IsExpired() bool {
	return time.Now().After(ce.Expires)
}
