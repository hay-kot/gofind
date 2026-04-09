package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/hay-kot/gofind/internal/gofind"
)

type CacheCmd struct {
	flags *Flags
}

func NewCacheCmd(flags *Flags) *CacheCmd {
	return &CacheCmd{flags: flags}
}

func (cmd *CacheCmd) Register(app *cli.Command) {
	app.Commands = append(app.Commands, &cli.Command{
		Name:   "cache",
		Usage:  "cache all config entries",
		Action: cmd.run,
	})
}

func (cmd *CacheCmd) run(ctx context.Context, c *cli.Command) error {
	cfg, err := readConfig(cmd.flags.ConfigFile)
	if err != nil {
		return err
	}

	finder := gofind.GoFind{Conf: cfg}
	start := time.Now()

	if err := finder.CacheAll(); err != nil {
		log.Err(err).Msg("failed to update cache")
		return err
	}

	log.Info().Dur("elapsed", time.Since(start)).Msg("cache updated")
	fmt.Println("cache updated in", time.Since(start))
	return nil
}
