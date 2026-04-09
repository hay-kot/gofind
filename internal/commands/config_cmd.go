package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/hay-kot/gofind/internal/paths"
	"github.com/hay-kot/gofind/internal/ui"
	"github.com/urfave/cli/v3"
)

type ConfigCmd struct {
	flags *Flags
}

func NewConfigCmd(flags *Flags) *ConfigCmd {
	return &ConfigCmd{flags: flags}
}

func (cmd *ConfigCmd) Register(app *cli.Command) {
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
}

func (cmd *ConfigCmd) list(ctx context.Context, c *cli.Command) error {
	p := paths.ConfigPath(cmd.flags.ConfigFile)
	cfg, err := readConfig(cmd.flags.ConfigFile)
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
