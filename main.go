package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hay-kot/gofind/internal/gofind"
	"github.com/hay-kot/gofind/internal/tui"
	"github.com/hay-kot/gofind/internal/ui"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Version: "0.3.0",
		Name:    "gofind",
		Usage:   "an interactive search for directories using the filepath.Match function",
		Commands: []*cli.Command{
			{
				Name:  "cache",
				Usage: "cache all config entries",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "entry",
						Usage: "specific entry to re-cache",
					},
				},
				Aliases: []string{"c"},
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg := gofind.ReadDefaultConfig()
					finder := gofind.GoFind{
						Conf: cfg,
					}

					start := time.Now()

					err := finder.CacheAll()
					if err != nil {
						log.Err(err).Msg("failed to update cache")
						return err
					}

					log.Info().Dur("elapsed", time.Since(start)).Msg("cache updated")
					return nil
				},
			},
			{
				Name:      "find",
				Usage:     "run interactive finder for entry",
				UsageText: "gofind find [config-entry string] e.g. `gofind find repos`",
				Flags:     []cli.Flag{},
				Aliases:   []string{"f"},
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg := gofind.ReadDefaultConfig()
					finder := gofind.GoFind{
						Conf: cfg,
					}

					entry := c.Args().Get(0)

					matches, err := finder.Run(entry)
					if err != nil {
						log.Err(err).Str("arg", entry).Msg("finder.Run failed")
						return err
					}

					result, err := tui.FuzzyFinder(matches)

					fmt.Println(result)
					return err
				},
			},
			{
				Name:  "setup",
				Usage: "first time setup",
				Action: func(ctx context.Context, c *cli.Command) error {
					if _, err := os.Stat(gofind.DefaultConfigPath()); err == nil {
						log.Warn().Str("path", gofind.DefaultConfigPath()).Msg("config file already exists")
						return nil
					}

					err := gofind.ConfigSetup()

					if err != nil {
						log.Err(err).Msg("config setup failed")
						return err
					}

					log.Info().Str("path", gofind.DefaultConfigPath()).Msg("config file created")
					return nil
				},
			},
			{
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
						Action: func(ctx context.Context, c *cli.Command) error {
							cfg := gofind.ReadDefaultConfig()
							p := gofind.DefaultConfigPath()

							if c.Bool("path") {
								fmt.Println(p)
								return nil
							}

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

							items := [][]string{
								{"Arg", "Root", "Match"},
							}

							for key, entry := range cfg.Commands {
								items = append(items, []string{key, strings.Join(entry.Roots, ", "), entry.MatchStr})
							}

							str.WriteString(ui.Table(items))
							fmt.Println(str.String())
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("gofind failed")
	}
}
