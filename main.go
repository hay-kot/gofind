package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/gofind/internal/core/config"
	"github.com/hay-kot/gofind/internal/gofind"
	"github.com/hay-kot/gofind/internal/tui"
	"github.com/hay-kot/gofind/internal/ui"
	"github.com/muesli/termenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func main() {
	// Sets colors so they show up in stderr:
	//
	// Stolen from
	//
	// - https://github.com/charmbracelet/gum/blob/6d405c49b1b929b771cda5fe939cf5900e392b70/main.go#L31
	lipgloss.SetColorProfile(termenv.NewOutput(os.Stderr).Profile)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().Caller().Logger().
		Level(zerolog.WarnLevel)

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
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg, err := config.ReadFile(config.XDGConfigPath())
					if err != nil {
						return err
					}

					finder := gofind.GoFind{Conf: cfg}

					start := time.Now()

					err = finder.CacheAll()
					if err != nil {
						log.Err(err).Msg("failed to update cache")
						return err
					}

					log.Info().Dur("elapsed", time.Since(start)).Msg("cache updated")
					fmt.Println("cached updated in", time.Since(start))
					return nil
				},
			},
			{
				Name:      "find",
				Usage:     "run interactive finder for entry",
				UsageText: "gofind find [config-entry string] e.g. `gofind find repos`",
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg, err := config.ReadFile(config.XDGConfigPath())
					if err != nil {
						return err
					}

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
							cfg, err := config.ReadFile(config.XDGConfigPath())
							if err != nil {
								return err
							}

							p := config.XDGConfigPath()

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
