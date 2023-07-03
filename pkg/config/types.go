package config

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v3"
)

func NewCLIFromJson(j []byte) (*CLI, error) {
	c := &CLI{
		FlagDelim:                " ",
		DoubleDashesForLongFlags: true,
	}

	err := json.Unmarshal(j, c)
	return c, err
}

func NewCLIFromYaml(y []byte) (*CLI, error) {
	c := &CLI{
		FlagDelim:                " ",
		DoubleDashesForLongFlags: true,
	}

	err := yaml.Unmarshal(y, c)
	return c, err
}

// tags generated using: gomodifytags -file pkg/config/types.go -all -add-tags json,yaml -transform camelcase
// Ref: https://github.com/fatih/gomodifytags
type CLI struct {
	Name                     string  `json:"name" yaml:"name"` // an arbitrary for the generated UI
	Help                     string  `json:"help,omitempty" yaml:"help,omitempty"`
	FlagDelim                string  `json:"flagDelim,omitempty" yaml:"flagDelim,omitempty"` // the delimiter used for flags between key and value (e.g. FlagDelim="=" will have --key=value)
	DoubleDashesForLongFlags bool    `json:"doubleDashesForLongFlags,omitempty" yaml:"doubleDashesForLongFlags,omitempty"`
	Command                  Command `json:"command" yaml:"command"` // the entry of the CLI, make sure the name to this Command is the path to the binary to be called
}

type Command struct {
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Flags       []FlagOrArg `json:"flags,omitempty" yaml:"flags,omitempty"`
	Args        []FlagOrArg `json:"args,omitempty" yaml:"args,omitempty"`
	Subcommands []Command   `json:"subcommands,omitempty" yaml:"subcommands,omitempty"` // e.g. kubectl get <resource>, here get is a subcommand to kubectl
}

type FlagArgType string

const (
	FlagArgTypeString  FlagArgType = "string"
	FlagArgTypeNumber  FlagArgType = "number"
	FlagArgTypeArray   FlagArgType = "array"
	FlagArgTypeBoolean FlagArgType = "boolean"
	FlagArgTypeEnum    FlagArgType = "enum"
)

type FlagOrArg struct {
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Type        FlagArgType `json:"type" yaml:"type"`
	Required    bool        `json:"required,omitempty" yaml:"required,omitempty"`
	Default     string      `json:"default,omitempty" yaml:"default,omitempty"`
	Options     []string    `json:"options,omitempty" yaml:"options,omitempty"` // only required when Type=enum
}
