package config_test

import (
	"testing"
)

func TestScriptGeneration(t *testing.T) {
	f := docker.Form()

	f.Flags["config"].Value = "this-config.yaml"
	f.Flags["config"].Long = true

	f.Flags["log-level"].Value = "info"
	f.Flags["log-level"].Long = true
	f.Choice = "volume"

	f.Subcommands["volume"].Choice = "create"
	f.Subcommands["volume"].Subcommands["create"].Args["name"].Value = "new-volume"

	s := docker.Script(f)

	if s != "docker --config this-config.yaml --log-level info volume create new-volume" {
		t.Errorf("unexpected script generated: %s", s)
	}
}
