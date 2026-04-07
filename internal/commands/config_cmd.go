package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/hay-kot/gofind/internal/config"
	"github.com/hay-kot/gofind/internal/ui"
)

// ConfigCmd implements the config subcommand.
type ConfigCmd struct {
	flags *Flags
}

// NewConfigCmd creates a new config command.
func NewConfigCmd(flags *Flags) *ConfigCmd {
	return &ConfigCmd{flags: flags}
}

// Register adds the config command to the application.
func (cmd *ConfigCmd) Register(app *cli.Command) *cli.Command {
	app.Commands = append(app.Commands, &cli.Command{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "add, remove, or list configuration entries",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list all config entries",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "path",
						Usage: "returns only the path",
					},
				},
				Action: cmd.list,
			},
		},
	})
	return app
}

func (cmd *ConfigCmd) list(ctx context.Context, c *cli.Command) error {
	p := config.XDGConfigPath(cmd.flags.ConfigFile)
	cfg, err := config.ReadFile(p)
	if err != nil {
		return err
	}

	if c.Bool("path") {
		fmt.Println(p)
		return nil
	}

	ui.Init(cfg.Theme)

	str := strings.Builder{}
	str.WriteString("\n")
	str.WriteString(ui.Bold("Config Path: ") + p + "\n\n")

	values := [][]string{
		{"Key", "Option", "Value"},
		{"default", "Default Argument", cfg.Default},
		{"cache", "Cache Dir", cfg.CacheDir},
		{"ignore", "Ignore Patterns", strings.Join(cfg.Ignore, ", ")},
		{"max_recursion", "Max Recursion", fmt.Sprintf("%d", cfg.MaxRecursion)},
	}

	str.WriteString(ui.Table(values))
	str.WriteString("\n")

	items := [][]string{{"Arg", "Root", "Match"}}
	for key, entry := range cfg.Commands {
		items = append(items, []string{key, strings.Join(entry.Roots, ", "), entry.MatchStr})
	}

	str.WriteString(ui.Table(items))
	fmt.Println(str.String())
	return nil
}
