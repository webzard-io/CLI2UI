package main

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
	"encoding/json"
	"log"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			if c.Args().Len() != 1 {
				return cli.Exit("Usage: cli2ui <config>", 1)
			}

			f, err := os.ReadFile(c.Args().First())
			if err != nil {
				return err
			}

			cfg, err := config.NewCLI(f)
			if err != nil {
				return err
			}

			ui, err := ui.NewUI(*cfg)
			if err != nil {
				return err
			}

			return ui.Run()
		},
		Commands: []*cli.Command{
			{
				Name: "jsonschema",
				Action: func(c *cli.Context) error {
					schema := jsonschema.Reflect(&config.CLI{})
					s, _ := json.MarshalIndent(schema, "", "    ")
					return os.WriteFile("cli.schema.json", s, 0644)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
