package paths

import (
	"os"
	"path/filepath"
)

const appName = "gofind"

// ConfigDir returns $XDG_CONFIG_HOME/gofind or ~/.config/gofind.
func ConfigDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}

	home, _ := homeDir()
	return filepath.Join(home, ".config", appName)
}

// ConfigPath resolves the config file path.
// Returns override if non-empty, otherwise ConfigDir()/gofind.json.
func ConfigPath(override string) string {
	if override != "" {
		return override
	}
	return filepath.Join(ConfigDir(), "gofind.json")
}

// DataDir returns $XDG_DATA_HOME/gofind or ~/.local/share/gofind.
func DataDir() string {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		return filepath.Join(dir, appName)
	}

	home, _ := homeDir()
	return filepath.Join(home, ".local", "share", appName)
}

// homeDir returns the user home directory, falling back to "/" on error.
func homeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "/", err
	}
	return home, nil
}
