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
	Theme        Theme                  `json:"theme,omitempty"`
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
