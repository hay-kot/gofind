package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/hay-kot/gofind/internal/paths"
)

type Theme struct {
	Prompt              string `json:"prompt"`
	SelectedIndicator   string `json:"selected_indicator"`
	SecondaryText       string `json:"secondary_text"`
	PrimaryText         string `json:"primary_text"`
	SelectionBackground string `json:"selection_background"`
}

type Config struct {
	Default      string                 `json:"default"`
	Commands     map[string]SearchEntry `json:"commands"`
	CacheDir     string                 `json:"cache"`
	Ignore       []string               `json:"ignore"`
	MaxRecursion int                    `json:"max_recursion"`
	Theme        Theme                  `json:"theme,omitzero"`
}

func DefaultTheme() Theme {
	return Theme{
		Prompt:              "#255F85",
		SelectedIndicator:   "#DA4167",
		SecondaryText:       "#848484",
		PrimaryText:         "#FFFFFF",
		SelectionBackground: "#2D2F27",
	}
}

func Default() *Config {
	return &Config{
		Default:      "",
		Commands:     make(map[string]SearchEntry),
		CacheDir:     filepath.Join(paths.DataDir(), "cache"),
		Ignore:       []string{},
		MaxRecursion: 10,
		Theme:        DefaultTheme(),
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

// XDGConfigPath resolves the config file path.
// Returns override if non-empty, otherwise returns ConfigDir()/gofind.json.
func XDGConfigPath(override string) string {
	if override != "" {
		return override
	}
	return filepath.Join(paths.ConfigDir(), "gofind.json")
}

// ReadFileOrDefault reads config from path, returning Default() if the file
// does not exist. Used in Before hooks where a missing config is acceptable.
func ReadFileOrDefault(path string) (*Config, error) {
	cfg, err := ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}
	return cfg, nil
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
	defer func() { _ = file.Close() }()

	return Write(file, cfg)
}
