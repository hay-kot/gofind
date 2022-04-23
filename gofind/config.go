package gofind

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Default      string                 `json:"default"`
	Commands     map[string]SearchEntry `json:"commands"`
	CacheDir     string                 `json:"cache"`
	Cache        map[string]CacheEntry  `json:"cache_old"`
	Ignore       []string               `json:"ignore"`
	MaxRecursion int                    `json:"max_recursion"`
}

func DefaultIgnore() []string {
	return []string{
		"node_modules",
		".venv",
		"venv",
	}
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

	if config.Ignore == nil {
		config.Ignore = DefaultIgnore()
	} else {
		config.Ignore = append(config.Ignore, DefaultIgnore()...)
	}

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

func (c Config) Save() {
	homedir := Must(os.UserHomeDir())
	configPath := filepath.Join(homedir, ".config", "gofind.json")

	file := Must(os.Create(configPath))
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	MustNotErr(encoder.Encode(c))
}
