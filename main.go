package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
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

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

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
		Version: build(),
		Name:    "gofind",
		Usage:   "an interactive search for directories using the filepath.Match function",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "specify configuration path",
				Sources:     cli.EnvVars("GOFIND_CONFIG"),
				Required:    false,
				Value:       "",
				DefaultText: config.XDGConfigPath(""),
				Aliases:     []string{"c"},
			},
			&cli.StringFlag{
				Name:   "cpuprofile",
				Usage:  "write cpu profile to file",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "memprofile",
				Usage:  "write memory profile to file",
				Hidden: true,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if cpuprofile := c.String("cpuprofile"); cpuprofile != "" {
				f, err := os.Create(cpuprofile)
				if err != nil {
					return ctx, fmt.Errorf("could not create CPU profile: %w", err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					_ = f.Close()
					return ctx, fmt.Errorf("could not start CPU profile: %w", err)
				}
				c.Metadata["cpuprofile_file"] = f
			}
			return ctx, nil
		},
		After: func(ctx context.Context, c *cli.Command) error {
			if f, ok := c.Metadata["cpuprofile_file"].(*os.File); ok {
				pprof.StopCPUProfile()
				_ = f.Close()
			}

			if memprofile := c.String("memprofile"); memprofile != "" {
				f, err := os.Create(memprofile)
				if err != nil {
					return fmt.Errorf("could not create memory profile: %w", err)
				}
				defer func() { _ = f.Close() }()
				runtime.GC()
				if err := pprof.WriteHeapProfile(f); err != nil {
					return fmt.Errorf("could not write memory profile: %w", err)
				}
			}
			return nil
		},
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
					cfg, err := config.ReadFile(config.XDGConfigPath(c.String("config")))
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
					cfg, err := config.ReadFile(config.XDGConfigPath(c.String("config")))
					if err != nil {
						return err
					}

					// Initialize theme from config
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
							p := config.XDGConfigPath(c.String("config"))
							cfg, err := config.ReadFile(p)
							if err != nil {
								return err
							}

							if c.Bool("path") {
								fmt.Println(p)
								return nil
							}

							// Initialize theme from config
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
