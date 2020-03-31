package app

import (
	"os"

	"github.com/rs/zerolog/log"
	cli "github.com/urfave/cli/v2"
)

func Execute() {
	app := &cli.App{
		Name:  "Question service",
		Usage: "RESTful micro-service application",

		Commands: []*cli.Command{
			{
				Action:  run,
				Name:    "runserver",
				Usage:   `launch server`,
				Aliases: []string{"r"},
				Flags:   runFlags,
			},
		},

		Action: run,
		Flags:  runFlags,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("start cli")
	}
}
