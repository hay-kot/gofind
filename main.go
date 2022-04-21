package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hay-kot/gofind/gofind"
	"github.com/urfave/cli/v2"
)

func main() {
	finder := gofind.App{}

	app := &cli.App{
		Name:      "gofind",
		Usage:     "an interactive search for directories using the filepath.Match function",
		UsageText: "gofind [config-entry string] e.g. `gofind repos`",
		Action: func(c *cli.Context) error {
			entry := c.Args().Get(0)
			result := finder.Run(entry)
			fmt.Println(result)
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
				Aliases: []string{"c"},
				Action: func(c *cli.Context) error {
					start := time.Now()

					err := finder.CacheAll()
					if err != nil {
						return err
					}

					fmt.Println("Caches Updated In:", time.Since(start))
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
