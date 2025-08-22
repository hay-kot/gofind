package gofind

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "test-cache")

	cache, err := NewCache(testDir)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	if cache.Dir != testDir {
		t.Errorf("Expected Dir to be %q, got %q", testDir, cache.Dir)
	}

	// Check that directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Expected cache directory to be created, but it doesn't exist")
	}
}

func TestNewCacheWithTildePath(t *testing.T) {
	cache, err := NewCache("~/test-cache")
	if err != nil {
		t.Fatalf("NewCache with tilde failed: %v", err)
	}

	if cache.Dir != "~/test-cache" {
		t.Errorf("Expected Dir to be ~/test-cache, got %q", cache.Dir)
	}
}

func TestCachePathConstructor(t *testing.T) {
	tmpDir := t.TempDir()
	cache := Cache{Dir: tmpDir}

	tests := []struct {
		namespace string
		expected  string
	}{
		{
			namespace: "test",
			expected:  filepath.Join(tmpDir, "test.json"),
		},
		{
			namespace: "repos",
			expected:  filepath.Join(tmpDir, "repos.json"),
		},
		{
			namespace: "projects-with-dashes",
			expected:  filepath.Join(tmpDir, "projects-with-dashes.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.namespace, func(t *testing.T) {
			result, err := cache.PathConstructor(tt.namespace)
			if err != nil {
				t.Fatalf("PathConstructor failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestCacheSetAndFind(t *testing.T) {
	tmpDir := t.TempDir()
	cache, err := NewCache(tmpDir)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	// Test data
	matches := []Match{
		{Name: "project1", Path: "/home/user/project1"},
		{Name: "project2", Path: "/home/user/project2"},
	}

	// Test Set
	entry, err := cache.Set("test-namespace", matches)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	if entry == nil {
		t.Fatal("Set returned nil entry")
	}

	if len(entry.Matches) != len(matches) {
		t.Errorf("Expected %d matches, got %d", len(matches), len(entry.Matches))
	}

	for i, match := range matches {
		if i < len(entry.Matches) {
			if entry.Matches[i].Name != match.Name {
				t.Errorf("Expected match name %q, got %q", match.Name, entry.Matches[i].Name)
			}
			if entry.Matches[i].Path != match.Path {
				t.Errorf("Expected match path %q, got %q", match.Path, entry.Matches[i].Path)
			}
		}
	}

	// Test Find
	foundEntry, err := cache.Find("test-namespace")
	if err != nil {
		t.Fatalf("Find failed: %v", err)
	}

	if foundEntry == nil {
		t.Fatal("Find returned nil entry")
	}

	if len(foundEntry.Matches) != len(matches) {
		t.Errorf("Expected %d matches in found entry, got %d", len(matches), len(foundEntry.Matches))
	}

	// Verify the data matches
	for i, expectedMatch := range matches {
		if i < len(foundEntry.Matches) {
			actualMatch := foundEntry.Matches[i]
			if actualMatch.Name != expectedMatch.Name {
				t.Errorf("Expected found match name %q, got %q", expectedMatch.Name, actualMatch.Name)
			}
			if actualMatch.Path != expectedMatch.Path {
				t.Errorf("Expected found match path %q, got %q", expectedMatch.Path, actualMatch.Path)
			}
		}
	}
}

func TestCacheFindNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	cache, err := NewCache(tmpDir)
	if err != nil {
		t.Fatalf("NewCache failed: %v", err)
	}

	entry, err := cache.Find("nonexistent-namespace")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound, got %v", err)
	}

	if entry != nil {
		t.Error("Expected nil entry for nonexistent cache")
	}
}

func TestCacheEntryExpiration(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		expires  time.Time
		expected bool
	}{
		{
			name:     "not expired",
			expires:  now.Add(1 * time.Hour),
			expected: false,
		},
		{
			name:     "expired",
			expires:  now.Add(-1 * time.Hour),
			expected: true,
		},
		{
			name:     "just expired",
			expires:  now.Add(-1 * time.Second),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := CacheEntry{
				Matches: []Match{},
				Expires: tt.expires,
			}

			if entry.IsExpired() != tt.expected {
				t.Errorf("Expected IsExpired() to return %v, got %v", tt.expected, entry.IsExpired())
			}
		})
	}
}

func TestCacheFindExpiredEntry(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an expired cache entry manually
	expiredEntry := CacheEntry{
		Matches: []Match{{Name: "test", Path: "/test"}},
		Expires: time.Now().Add(-1 * time.Hour), // Expired 1 hour ago
	}

	cachePath := filepath.Join(tmpDir, "expired.json")
	file, err := os.Create(cachePath)
	if err != nil {
		t.Fatalf("Failed to create cache file: %v", err)
	}

	encoder := json.NewEncoder(file)
	err = encoder.Encode(expiredEntry)
	if err != nil {
		_ = file.Close()
		t.Fatalf("Failed to encode expired entry: %v", err)
	}
	_ = file.Close()

	cache := Cache{Dir: tmpDir}
	
	// Try to find the expired entry
	entry, err := cache.Find("expired")
	if err != ErrCacheNotFound {
		t.Errorf("Expected ErrCacheNotFound for expired cache, got %v", err)
	}

	if entry != nil {
		t.Error("Expected nil entry for expired cache")
	}
}

func TestCacheSetCreatesDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Use nested path that doesn't exist
	nestedCacheDir := filepath.Join(tmpDir, "nested", "cache", "dir")
	cache := Cache{Dir: nestedCacheDir}

	matches := []Match{{Name: "test", Path: "/test"}}
	
	_, err := cache.Set("test", matches)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify the nested directories were created
	if _, err := os.Stat(nestedCacheDir); os.IsNotExist(err) {
		t.Error("Expected nested cache directory to be created")
	}

	// Verify the cache file exists
	cachePath := filepath.Join(nestedCacheDir, "test.json")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Error("Expected cache file to be created")
	}
}