package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/hay-kot/gofind/internal/config"
	"github.com/hay-kot/gofind/internal/gofind"
)

// CacheCmd implements the cache subcommand.
type CacheCmd struct {
	flags *Flags
}

// NewCacheCmd creates a new cache command.
func NewCacheCmd(flags *Flags) *CacheCmd {
	return &CacheCmd{flags: flags}
}

// Register adds the cache command to the application.
func (cmd *CacheCmd) Register(app *cli.Command) *cli.Command {
	app.Commands = append(app.Commands, &cli.Command{
		Name:  "cache",
		Usage: "cache all config entries",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "entry",
				Usage: "specific entry to re-cache",
			},
		},
		Action: cmd.run,
	})
	return app
}

func (cmd *CacheCmd) run(ctx context.Context, c *cli.Command) error {
	cfg, err := config.ReadFile(config.XDGConfigPath(cmd.flags.ConfigFile))
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
