package config_test

import (
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

	if s != "docker --config this-config.yaml --log-level info volume create new-volume" &&
		s != "docker --log-level info --config this-config.yaml volume create new-volume" {
		t.Errorf("unexpected script generated: %s", s)
	}
}
