package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Default      string                 `json:"default"`
	Commands     map[string]SearchEntry `json:"commands"`
	CacheDir     string                 `json:"cache"`
	Ignore       []string               `json:"ignore"`
	MaxRecursion int                    `json:"max_recursion"`
}

func Default() *Config {
	return &Config{
		Default:      "",
		Commands:     make(map[string]SearchEntry),
		CacheDir:     XDGCachePath(),
		Ignore:       []string{},
		MaxRecursion: 10,
	}
}

// IgnorePatterns are common ignore folders/file patterns that
// will always be excluded.
func IgnorePatterns() []string {
	return []string{
		"node_modules",
		".venv",
		"venv",
	}
}

type SearchEntry struct {
	Roots    []string `json:"roots"`
	MatchStr string   `json:"match"`
}

func Read(r io.Reader) (*Config, error) {
	cfg := Default()
	err := json.NewDecoder(r).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func ReadFile(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return Read(bytes.NewBuffer(file))
}

func Write(w io.Writer, cfg *Config) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}

func WriteFile(path string, cfg *Config) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return Write(file, cfg)
}
