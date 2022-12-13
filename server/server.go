package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/y1rn/nsync/neko"
	"github.com/y1rn/nsync/sync"
	"gopkg.in/yaml.v3"
)

type Config struct {
	NodePrefix     string
	Token          string
	ConfigFilePath string
	Url            string
	Address        string
}

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
			&cli.StringFlag{
				Name:        "address",
				Aliases:     []string{"a"},
				Destination: &config.Address,
				Value:       ":8080",
				Usage:       "server address",
			},
		},
		Action: func(_ *cli.Context) error {
			configs, err := sync.LoadConfig(config.ConfigFilePath)
			if err != nil {
				return err
			}

			for _, c := range configs {
				http.HandleFunc(c.OutPut, func(c sync.Config) func(w http.ResponseWriter, r *http.Request) {
					return func(w http.ResponseWriter, _ *http.Request) {
						ruleResp, err := neko.LoadRule(config.Url, config.Token)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(err.Error()))
						}
						if len(ruleResp.Data) == 0 {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(errors.New("no neko rule data").Error()))
						}
						temp, err := sync.GetOutput(c, ruleResp.Data, config.NodePrefix)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(err.Error()))
						}
						en := yaml.NewEncoder(w)
						en.SetIndent(1)
						err = en.Encode(temp)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(err.Error()))
						}

					}
				}(c))
			}
			http.ListenAndServe(config.Address, nil)

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		// os.Exit(1)
		log.Fatal(err)
	}
}
