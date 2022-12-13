package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/y1rn/nsync/sync"
)

var config Config

func main() {
	app := &cli.App{
		Name:  "nsync",
		Usage: "sync nekoneko nodes to your config file",
		// EnableBashCompletion: true,
		// Suggest:              true,
		HideHelp: true,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "prefix",
				Aliases:     []string{"p"},
				Destination: &config.NodePrefix,
				Value:       "neko",
				Usage:       "node name prefix to config files, example: neko-ipel",
			},
			&cli.PathFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Required:    true,
				Destination: &config.ConfigFilePath,
				Usage:       "yaml config file path",
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Required:    true,
				Destination: &config.Token,
				Usage:       "nekoneko api token",
			},
			&cli.StringFlag{
				Name:        "url",
				Destination: &config.Url,
				Value:       "https://relay.nekoneko.cloud/api/rules",
				Usage:       "nekoneko api url",
			},
		},
		Action: func(ctx *cli.Context) error {
			return sync.Sync(config.Url, config.Token, config.ConfigFilePath, config.NodePrefix)
		},
	}

	if err := app.Run(os.Args); err != nil {
		// os.Exit(1)
		log.Fatal(err)
	}
}
