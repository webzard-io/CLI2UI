package config

type FlagArgType string

const (
	FlagArgTypeString  FlagArgType = "string"
	FlagArgTypeNumber  FlagArgType = "number"
	FlagArgTypeArray   FlagArgType = "array"
	FlagArgTypeBoolean FlagArgType = "boolean"
	FlagArgTypeEnum    FlagArgType = "enum"
)

type CLI struct {
	Name    string  `json:"name" yaml:"name"`           // CLI name, e.g., "ping", "kubectl", etc.
	Help    string  `json:"help,omitempty" yaml:"help"` // CLI help message
	Command Command `json:"command" yaml:"command"`     // List of supported actions by the CLI

	FlagDelim string `json:",omitempty" yaml:",omitempty"`
}

type FlagOrArg struct {
	Name        string      `json:"name" yaml:"name"` // Name of the flag parameter
	Description string      `json:"description,omitempty" yaml:"description"`
	Type        FlagArgType `json:"type,omitempty" yaml:"type,omitempty"`         // Type of the flag parameter, optional field
	Required    bool        `json:"required,omitempty" yaml:"required,omitempty"` // Whether the flag parameter is required, optional field
	Default     string      `json:"default,omitempty" yaml:"default,omitempty"`   // Default value for the flag parameter, optional field
	Options     []string    `json:"options,omitempty" yaml:"options,omitempty"`   // Enum options when flag type is enum
}

type Command struct {
	Name        string      `json:"name" yaml:"name"` // Name of the action, e.g., "apply", "delete", etc.
	Description string      `json:"description,omitempty" yaml:"description"`
	Flags       []FlagOrArg `json:"flags,omitempty" yaml:"flags"`             // List of flag parameters for the action
	Args        []FlagOrArg `json:"args,omitempty" yaml:"args,omitempty"`     // List of positional arguments for the action, optional field
	Subcommands []Command   `json:"subcommands,omitempty" yaml:"subcommands"` // List of supported actions by the CLI
}
