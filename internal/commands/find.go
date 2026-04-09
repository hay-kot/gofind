package commands

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/hay-kot/gofind/internal/gofind"
	"github.com/hay-kot/gofind/internal/tui"
	"github.com/hay-kot/gofind/internal/ui"
)

type FindCmd struct {
	flags *Flags
}

func NewFindCmd(flags *Flags) *FindCmd {
	return &FindCmd{flags: flags}
}

func (cmd *FindCmd) Register(app *cli.Command) {
	var noCache bool
	app.Commands = append(app.Commands, &cli.Command{
		Name:      "find",
		Usage:     "run interactive finder for entry",
		UsageText: "gofind find [config-entry string] e.g. `gofind find repos`",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "no-cache",
				Usage:       "skip cache and scan directories directly",
				Destination: &noCache,
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			return cmd.run(ctx, c, noCache)
		},
	})
}

func (cmd *FindCmd) run(ctx context.Context, c *cli.Command, noCache bool) error {
	cfg, err := readConfig(cmd.flags.ConfigFile)
	if err != nil {
		return err
	}

	ui.Init(cfg.Theme)

	finder := gofind.GoFind{Conf: cfg, NoCache: noCache}
	entry := c.Args().Get(0)

	matches, err := finder.Run(entry)
	if err != nil {
		log.Err(err).Str("arg", entry).Msg("finder.Run failed")
		return err
	}

	result, err := tui.FuzzyFinder(matches)
	if result != "" {
		fmt.Println(result)
	}
	return err
}
