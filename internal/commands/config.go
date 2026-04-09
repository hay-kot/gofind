package commands

import (
	"github.com/hay-kot/gofind/internal/config"
	"github.com/hay-kot/gofind/internal/paths"
)

func readConfig(override string) (*config.Config, error) {
	return config.ReadFile(paths.ConfigPath(override))
}
