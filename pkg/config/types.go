package config

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v3"
)

func NewCLIFromJson(j []byte) (*CLI, error) {
	c := &CLI{
		FlagDelim: " ",
	}

	err := json.Unmarshal(j, c)
	return c, err
}

func NewCLIFromYaml(y []byte) (*CLI, error) {
	c := &CLI{
		FlagDelim: " ",
	}

	err := yaml.Unmarshal(y, c)
	return c, err
}

// tags generated using: gomodifytags -file pkg/config/types.go -all -add-tags json,yaml -transform camelcase
// Ref: https://github.com/fatih/gomodifytags
type CLI struct {
	Name      string  `json:"name" yaml:"name"` // an arbitrary for the generated UI
	Help      string  `json:"help" yaml:"help"`
	FlagDelim string  `json:"flagDelim" yaml:"flagDelim"` // the delimiter used for flags between key and value (e.g. FlagDelim="=" will have --key=value)
	Command   Command `json:"command" yaml:"command"`     // the entry of the CLI, make sure the name to this Command is the path to the binary to be called
}

type Command struct {
	Name        string      `json:"name" yaml:"name"`
	Description string      `json:"description" yaml:"description"`
	Flags       []FlagOrArg `json:"flags" yaml:"flags"`
	Args        []FlagOrArg `json:"args" yaml:"args"`
	Subcommands []Command   `json:"subcommands" yaml:"subcommands"` // e.g. kubectl get <resource>, here get is a subcommand to kubectl
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
	Description string      `json:"description" yaml:"description"`
	Type        FlagArgType `json:"type" yaml:"type"`
	Required    bool        `json:"required" yaml:"required"`
	Default     string      `json:"default" yaml:"default"`
	Options     []string    `json:"options" yaml:"options"` // only required when Type=enum
}
