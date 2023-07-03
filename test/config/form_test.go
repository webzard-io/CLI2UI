package config_test

import (
	"CLI2UI/pkg/config"
	"testing"
)

func TestScriptGeneration(t *testing.T) {
	f := docker.Form()

	f.Flags["config"] = "this-config.yaml"
	f.Flags["log-level"] = "info"
	f.Choice = "volume"

	f.Subcommands["volume"].Choice = "create"
	f.Subcommands["volume"].Subcommands["create"].Args["name"] = "new-volume"

	s := docker.Script(f)

	if s != "docker --config this-config.yaml --log-level info volume create new-volume" {
		t.Errorf("unexpected script generated: %s", s)
	}
}

var docker = config.CLI{
	Name:      "docker",
	FlagDelim: " ",
	Command: config.Command{
		Name: "docker",
		Flags: []config.FlagOrArg{
			{
				Name: "config",
				Type: config.FlagArgTypeString,
			},
			{
				Name:    "log-level",
				Type:    config.FlagArgTypeEnum,
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
