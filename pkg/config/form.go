package config

import "fmt"

type Form struct {
	Flags       map[string]*FlagOrArgValue
	Args        map[string]*FlagOrArgValue
	Subcommands map[string]*Form
	Choice      string
}

type FlagOrArgValue struct {
	Value any
	Long  bool
}

func (c CLI) Script(f Form) string {
	return parseScript(&f, c.Name, c.FlagDelim)
}

func parseScript(f *Form, script string, flagDelim string) string {
	for k, v := range f.Flags {
		if v.Value == nil {
			continue
		}

		prefix := "-"

		if v.Long {
			prefix += "-"
		}

		script = fmt.Sprintf("%s %s%s%s%s", script, prefix, k, flagDelim, v.Value)
	}

	for _, v := range f.Args {
		if v.Value == nil {
			continue
		}
		script = fmt.Sprintf("%s %s", script, v.Value)
	}

	if len(f.Subcommands) == 0 || f.Choice == "" {
		return script
	}

	script = fmt.Sprintf("%s %s", script, f.Choice)
	return parseScript((f.Subcommands)[f.Choice], script, flagDelim)
}

func (c CLI) Form() Form {
	f := Form{}

	parseForm(&c.Command, &f)

	return f
}

func parseForm(c *Command, f *Form) {
	if c == nil {
		return
	}

	flags := map[string]*FlagOrArgValue{}
	args := map[string]*FlagOrArgValue{}
	subcommands := map[string]*Form{}

	f.Flags = flags
	f.Args = args
	f.Subcommands = subcommands

	for _, f := range c.Flags {
		flags[f.Name] = &FlagOrArgValue{
			Value: nil,
			Long:  f.Long,
		}
	}

	for _, a := range c.Args {
		args[a.Name] = &FlagOrArgValue{
			Value: nil,
		}
	}

	for _, c := range c.Subcommands {
		form := &Form{}
		parseForm(&c, form)
		f.Subcommands[c.Name] = form
	}
}
