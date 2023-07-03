package config_test

import (
	"CLI2UI/pkg/config"
	"os"
	"reflect"
	"testing"
)

func TestNewCLIFromYaml(t *testing.T) {
	data, err := os.ReadFile("./docker.yaml")
	if err != nil {
		t.Error("failed reading docker.yaml")
	}

	constructed, err := config.NewCLIFromYaml(data)
	if err != nil {
		t.Error("failed building CLI from YAML")
	}

	equal := reflect.DeepEqual(docker, constructed)

	if !equal {
		t.Error("instance constructed by YAML does not equal to the declared one")
	}
}

func TestNewCLIFromJson(t *testing.T) {
	data, err := os.ReadFile("./docker.json")
	if err != nil {
		t.Error("failed reading docker.json")
	}

	constructed, err := config.NewCLIFromJson(data)
	if err != nil {
		t.Error("failed building CLI from JSON")
	}

	equal := reflect.DeepEqual(docker, constructed)

	if !equal {
		t.Error("instance constructed by JSON does not equal to the declared one")
	}
}

var docker = &config.CLI{
	Name:                     "docker",
	FlagDelim:                " ",
	DoubleDashesForLongFlags: true,
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
