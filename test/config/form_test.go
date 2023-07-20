package config_test

import (
	"testing"
)

func TestScriptGeneration(t *testing.T) {
	f := docker.Form()

	f.Flags["config"].Enabled = true
	f.Flags["config"].Value = "this-config.yaml"

	f.Flags["log-level"].Enabled = true
	f.Flags["log-level"].Value = "info"

	f.Choice = "volume"

	f.Subcommands["volume"].Choice = "create"
	f.Subcommands["volume"].Subcommands["create"].Args["name"].Enabled = true
	f.Subcommands["volume"].Subcommands["create"].Args["name"].Value = "new-volume"

	s, _ := docker.Script(f)

	if s != "docker --config this-config.yaml --log-level info volume create new-volume" {
		t.Errorf("unexpected script generated: %s", s)
	}
}
