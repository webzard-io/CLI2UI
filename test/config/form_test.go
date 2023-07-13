package config_test

import (
	"CLI2UI/pkg/config"
	"testing"
)

func TestScriptGeneration(t *testing.T) {
	f := docker.Form()

	f.Flags["config"] = &config.OptionValue{
		Value:   "this-config.yaml",
		Long:    true,
		Enabled: true,
	}

	f.Flags["log-level"] = &config.OptionValue{
		Value:   "info",
		Long:    true,
		Enabled: true,
	}

	f.Choice = "volume"

	f.Subcommands["volume"].Choice = "create"
	f.Subcommands["volume"].Subcommands["create"].Args["name"] = &config.OptionValue{
		Value:   "new-volume",
		Enabled: true,
	}

	s, _ := docker.Script(f)

	if s != "docker --config this-config.yaml --log-level info volume create new-volume" {
		t.Errorf("unexpected script generated: %s", s)
	}
}
