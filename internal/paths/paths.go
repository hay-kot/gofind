package paths

import (
	"os"
	"path/filepath"
)

const appName = "gofind"

// ConfigDir returns the XDG config directory for the application.
// Uses $XDG_CONFIG_HOME/<app> or falls back to ~/.config/<app>.
func ConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", appName)
}

// DataDir returns the XDG data directory for the application.
// Uses $XDG_DATA_HOME/<app> or falls back to ~/.local/share/<app>.
func DataDir() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", appName)
}

// CacheDir returns the XDG cache directory for the application.
// Uses $XDG_CACHE_HOME/<app> or falls back to ~/.cache/<app>.
func CacheDir() string {
	if dir := os.Getenv("XDG_CACHE_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", appName)
}
