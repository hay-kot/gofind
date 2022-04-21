package gofind

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	Default  string                 `json:"default"`
	Commands map[string]SearchEntry `json:"commands"`
	Cache    map[string]CacheEntry  `json:"cache"`
}

func DefaultConfigPath() string {
	homedir := Must(os.UserHomeDir())
	configPath := filepath.Join(homedir, ".config", "gofind.json")

	return configPath
}

func ReadConfig(path string) Config {
	config := Config{}

	file := Must(os.Open(path))

	decoder := json.NewDecoder(file)

	MustNotErr(decoder.Decode(&config))

	return config
}

func ReadDefaultConfig() Config {
	configPath := DefaultConfigPath()

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config file does not exist")
	}

	return ReadConfig(configPath)
}

func (c Config) CacheAll() {
	for key, search := range c.Commands {
		cache := search.Results()

		c.Cache[key] = CacheEntry{
			Matches: cache,
			Expires: time.Now().Add(time.Hour * 12),
		}
	}

	c.Save()
}

func (c Config) Save() {
	homedir := Must(os.UserHomeDir())
	configPath := filepath.Join(homedir, ".config", "gofind.json")

	file := Must(os.Create(configPath))
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	MustNotErr(encoder.Encode(c))
}
