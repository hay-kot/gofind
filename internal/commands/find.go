package commands

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/hay-kot/gofind/internal/config"
	"github.com/hay-kot/gofind/internal/gofind"
	"github.com/hay-kot/gofind/internal/tui"
	"github.com/hay-kot/gofind/internal/ui"
)

// FindCmd implements the find subcommand.
type FindCmd struct {
	flags *Flags
}

// NewFindCmd creates a new find command.
func NewFindCmd(flags *Flags) *FindCmd {
	return &FindCmd{flags: flags}
}

// Register adds the find command to the application.
func (cmd *FindCmd) Register(app *cli.Command) *cli.Command {
	app.Commands = append(app.Commands, &cli.Command{
		Name:      "find",
		Usage:     "run interactive finder for entry",
		UsageText: "gofind find [config-entry string] e.g. `gofind find repos`",
		Action:    cmd.run,
	})
	return app
}

func (cmd *FindCmd) run(ctx context.Context, c *cli.Command) error {
	cfg, err := config.ReadFile(config.XDGConfigPath(cmd.flags.ConfigFile))
	if err != nil {
		return err
	}

	ui.Init(cfg.Theme)

	finder := gofind.GoFind{Conf: cfg}
	entry := c.Args().Get(0)

	matches, err := finder.Run(entry)
	if err != nil {
		log.Err(err).Str("arg", entry).Msg("finder.Run failed")
		return err
	}

	result, err := tui.FuzzyFinder(matches)
	fmt.Println(result)
	return err
}
