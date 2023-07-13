package config

import (
	"encoding/json"
	"errors"

	yaml "gopkg.in/yaml.v3"
)

func NewCLI(b []byte) (*CLI, error) {
	c := &CLI{
		OptionDelim: " ",
	}

	err := json.Unmarshal(b, c)
	if err == nil {
		return c, nil
	}

	err = yaml.Unmarshal(b, c)
	if err == nil {
		return c, nil
	}

	return nil, errors.New("failed parsing CLI config")
}

func NewCLIFromJson(j []byte) (*CLI, error) {
	c := &CLI{
		OptionDelim: " ",
	}

	err := json.Unmarshal(j, c)
	return c, err
}

func NewCLIFromYaml(y []byte) (*CLI, error) {
	c := &CLI{
		OptionDelim: " ",
	}

	err := yaml.Unmarshal(y, c)
	return c, err
}

// tags generated using: gomodifytags -file pkg/config/types.go -all -add-tags json,yaml -transform camelcase
// Ref: https://github.com/fatih/gomodifytags
type CLI struct {
	Name         string  `json:"name" yaml:"name"` // an arbitrary for the generated UI
	Help         string  `json:"help,omitempty" yaml:"help,omitempty"`
	OptionDelim  string  `json:"optionDelim,omitempty" yaml:"optionDelim,omitempty"`   // the delimiter used for flags between key and value (e.g. FlagDelim="=" will have --key=value)
	Command      Command `json:"command" yaml:"command"`                               // the entry of the CLI, make sure the name to this Command is the path to the binary to be called
	ExplicitBool bool    `json:"explicitBool,omitempty" yaml:"explicitBool,omitempty"` // if true, boolean flags will be specified in the form of `--flag=true` instead of `--flag`
}

type Format string

const (
	UnknownFormat Format = ""
	FormatJson    Format = "json"
	FormatYaml    Format = "yaml"
)

type Command struct {
	Name        string    `json:"name" yaml:"name"`
	Display     string    `json:"display,omitempty" yaml:"display,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Flags       []Option  `json:"flags,omitempty" yaml:"flags,omitempty"`
	Args        []Option  `json:"args,omitempty" yaml:"args,omitempty"`
	Subcommands []Command `json:"subcommands,omitempty" yaml:"subcommands,omitempty"` // e.g. kubectl get <resource>, here get is a subcommand to kubectl
	Format      Format    `json:"format,omitempty" yaml:"format,omitempty"`
}

type OptionType string

const (
	OptionTypeString  OptionType = "string"
	OptionTypeNumber  OptionType = "number"
	OptionTypeArray   OptionType = "array"
	OptionTypeBoolean OptionType = "boolean"
	OptionTypeEnum    OptionType = "enum"
)

type Option struct {
	Name        string     `json:"name" yaml:"name"`
	Type        OptionType `json:"type" yaml:"type"`
	Display     string     `json:"display,omitempty" yaml:"display,omitempty"`
	Long        bool       `json:"long,omitempty" yaml:"long,omitempty"` // if true, the flag will be specified in the form of `--flag` instead of `-flag`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool       `json:"required,omitempty" yaml:"required,omitempty"`
	Default     any        `json:"default,omitempty" yaml:"default,omitempty"`
	Options     []string   `json:"options,omitempty" yaml:"options,omitempty"` // only required when Type=enum
}

func (o Option) DisplayName() string {
	name := o.Display
	if name == "" {
		name = o.Name
	}
	return name
}

func (c Command) DisplayName() string {
	name := c.Display
	if name == "" {
		name = c.Name
	}
	return name
}
