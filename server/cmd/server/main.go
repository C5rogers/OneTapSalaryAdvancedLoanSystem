package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "lms-server",
		Usage: "lms-server --config=config.yml",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Usage: "pass configuration path",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			if len(cliCtx.String("config")) == 0 {
				return run(nil)
			} else {
				path := cliCtx.String("config")
				return run(&path)
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
