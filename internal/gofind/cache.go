package gofind

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
)

var ErrCacheNotFound = errors.New("cache not found")

type Cache struct {
	Dir string
}

func (ce Cache) PathConstructor(namespace string) string {
	root := ParsePath(ce.Dir)

	return filepath.Join(root, fmt.Sprintf("%s.json", namespace))
}

func NewCache(dir string) Cache {
	p := ParsePath(dir)

	// Create directory if it doesn't exist
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Fatal().Err(err).Msgf("os.MkdirAll(p=%s) failed with error '%s'", p, err.Error())
		}
	}

	return Cache{
		Dir: dir,
	}
}

// Find will find the first match in the cache or return an error if the cache
// isn't found.
func (c Cache) Find(namespace string) (*CacheEntry, error) {
	cachePath := c.PathConstructor(namespace)

	// Check if file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, ErrCacheNotFound
	}

	// Load cache
	return nil, ErrCacheNotFound
}

func (c Cache) Set(namespace string, results []Match) (*CacheEntry, error) {
	p := c.PathConstructor(namespace)

	entry := CacheEntry{
		Matches: results,
		Expires: time.Now().Add(time.Hour * 24),
	}

	// Create Parent
	err := os.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return nil, err
	}

	file := Must(os.Create(p))
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
