package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hay-kot/gofind/gofind"
	"github.com/hay-kot/gofind/tui"
	"github.com/hay-kot/gofind/ui"
	"github.com/hay-kot/yal"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Version: "0.1.3",
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
				Action: func(c *cli.Context) error {
					cfg := gofind.ReadDefaultConfig()
					finder := gofind.GoFind{
						Conf: cfg,
					}

					start := time.Now()

					err := finder.CacheAll()
					if err != nil {
						yal.Error(err.Error())
						return err
					}

					yal.Infof("caches updated in: %s", time.Since(start))
					return nil
				},
			},
			{
				Name:      "find",
				Usage:     "run interactive finder for entry",
				UsageText: "gofind find [config-entry string] e.g. `gofind find repos`",
				Flags:     []cli.Flag{},
				Aliases:   []string{"f"},
				Action: func(c *cli.Context) error {
					cfg := gofind.ReadDefaultConfig()
					finder := gofind.GoFind{
						Conf: cfg,
					}

					entry := c.Args().Get(0)

					matches, err := finder.Run(entry)
					if err != nil {
						yal.Error("finder.Run(%s) failed: %s", entry, err.Error())
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
				Action: func(c *cli.Context) error {
					if _, err := os.Stat(gofind.DefaultConfigPath()); err == nil {
						yal.Errorf("config file already exists: %s", gofind.DefaultConfigPath())
						return nil
					}

					err := gofind.ConfigSetup()

					if err != nil {
						yal.Error(err.Error())
						return err
					}

					yal.Infof("config file created: %s", gofind.DefaultConfigPath())
					return nil
				},
			},
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "add, remove, or list configuration entries",
				Subcommands: []*cli.Command{
					{
						Name:  "add",
						Usage: "add a config entry",
						Action: func(c *cli.Context) error {
							cfg := gofind.ReadDefaultConfig()
							if c.NArg() < 3 {
								yal.Error("missing arguments")
								return nil
							}

							key := c.Args().Get(0)
							entry := gofind.SearchEntry{
								Roots:    []string{c.Args().Get(1)}, //TODO: Let user setup multiple roots!
								MatchStr: c.Args().Get(2),
							}
							cfg.Commands[key] = entry
							cfg.Save()

							yal.Infof("Key=%s, Root=%s, MatchStr=%s", key, entry.Roots, entry.MatchStr)
							yal.Info("config entry added successfully")
							return nil
						},
					},
					{
						Name:  "remove",
						Usage: "remove a config entry",
						Action: func(c *cli.Context) error {
							cfg := gofind.ReadDefaultConfig()
							if c.NArg() < 1 {
								yal.Error("missing arguments")
								return nil
							}

							key := c.Args().Get(0)
							delete(cfg.Commands, key)
							cfg.Save()

							yal.Infof("Key=%s removed successfully", key)
							return nil
						},
					},
					{
						Name:  "list",
						Usage: "list all config entries",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "path",
								Usage: "returns only the path",
							},
						},
						Action: func(c *cli.Context) error {
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
