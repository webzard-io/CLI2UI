package main

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
)

func main() {
	ui := ui.NewUI(*docker)
	err := ui.Run()
	if err != nil {
		panic(err)
	}
}

var docker = &config.CLI{
	Name:      "docker",
	FlagDelim: " ",
	Command: config.Command{
		Name: "docker",
		Flags: []config.FlagOrArg{
			{
				Name: "config",
				Type: config.FlagArgTypeString,
				Long: true,
			},
			{
				Name:    "log-level",
				Type:    config.FlagArgTypeEnum,
				Long:    true,
				Default: "info",
				Options: []string{"debug", "info", "warn", "error", "fatal"},
			},
		},
		Subcommands: []config.Command{
			{
				Name: "volume",
				Subcommands: []config.Command{
					{
						Name: "create",
						Flags: []config.FlagOrArg{
							{
								Name: "driver",
								Type: config.FlagArgTypeString,
								Long: true,
							},
						},
						Args: []config.FlagOrArg{
							{
								Name: "name",
								Type: config.FlagArgTypeString,
							},
						},
					},
				},
			},
		},
	},
}
