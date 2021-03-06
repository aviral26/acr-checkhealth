package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

// Set at linking time.
var (
	Version = "dev"
)

func main() {
	app := &cli.App{
		Name:    "acr",
		Usage:   "ACR Check Health - evaluate the health of a registry",
		Version: Version,
		Authors: []*cli.Author{
			{
				Name: "Aviral Takkar",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  traceStr,
				Usage: "print trace logs with secrets",
			},
		},
		Commands: []*cli.Command{
			pingCommand,
			checkHealthCommand,
			referrersCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}
