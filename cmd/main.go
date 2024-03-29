package main

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
	"CLI2UI/pkg/ui/flat"
	"CLI2UI/pkg/ui/naive"
	"encoding/json"
	"log"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "CLI2UI",
		Usage: "Usage: cli2ui <config>",
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

			var newUI func(config.CLI) (ui.UI, error)
			switch cfg.UI {
			case config.UINaive:
				newUI = naive.NewUI
			default:
				newUI = flat.NewUI
			}

			ui, err := newUI(*cfg)
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
