package config

import "encoding/json"

func NewCLIFromJson(j []byte) (*CLI, error) {
	c := &CLI{
		// default value
		FlagDelim: " ",
	}

	err := json.Unmarshal(j, c)
	return c, err
}

type CLI struct {
	Name      string // an arbitrary for the generated UI
	Help      string
	FlagDelim string  // the delimiter used for flags between key and value (e.g. FlagDelim="=" will have --key=value)
	Command   Command // the entry of the CLI, make sure the name to this Command is the path to the binary to be called
}

type Command struct {
	Name        string
	Description string
	Flags       []FlagOrArg
	Args        []FlagOrArg
	Subcommands []Command // e.g. kubectl get <resource>, here get is a subcommand to kubectl
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
	Name        string
	Description string
	Type        FlagArgType
	Required    bool
	Default     string
	Options     []string // only required when Type=enum
}
