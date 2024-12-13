package config

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

const AppName = "gofind"
const ConfigFile = AppName + ".json"

func XDGConfigPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		configDir = filepath.Join(home, ".config")
	}

	return filepath.Join(configDir, AppName, ConfigFile)
}

func XDGCachePath() string {
	cacheDir := os.Getenv("XDG_DATA_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get home directory")
		}

		cacheDir = filepath.Join(home, ".local", "share")
	}

	return filepath.Join(cacheDir, AppName, "cache")
}
