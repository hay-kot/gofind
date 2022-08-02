package gofind

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hay-kot/yal"
)

type Config struct {
	Default      string                 `json:"default"`
	Commands     map[string]SearchEntry `json:"commands"`
	CacheDir     string                 `json:"cache"`
	Ignore       []string               `json:"ignore"`
	MaxRecursion int                    `json:"max_recursion"`
}

type SearchEntry struct {
	Root     string `json:"root"`
	MatchStr string `json:"match"`
}

func ConfigSetup() error {
	path := DefaultConfigPath()
	if _, err := os.Stat(path); os.IsExist(err) {
		yal.Errorf("config file already exists %s", path)
		os.Exit(1)
	}

	c := Config{
		Default:      "",
		Commands:     make(map[string]SearchEntry),
		CacheDir:     "",
		Ignore:       []string{},
		MaxRecursion: 10,
	}

	c.Save()
	return nil
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

	NoErr(decoder.Decode(&config))

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
		yal.Errorf("Config file not found %s", configPath)
		os.Exit(1)
	}

	return ReadConfig(configPath)
}

func (c Config) Save() {
	homedir := Must(os.UserHomeDir())
	configPath := filepath.Join(homedir, ".config", "gofind.json")

	file := Must(os.Create(configPath))
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	NoErr(encoder.Encode(c))
}
