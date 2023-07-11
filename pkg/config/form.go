package config

import "fmt"

type Form struct {
	Flags       map[string]*OptionValue
	Args        map[string]*OptionValue
	Subcommands map[string]*Form
	Choice      string
}

type OptionValue struct {
	Value any
	Long  bool
}

func (c CLI) Script(f Form) string {
	return parseScript(&f, c.Command.Name, c.OptionDelim, c.ExplicitBool)
}

func parseScript(f *Form, script string, optionDelim string, explicitBool bool) string {
	for k, v := range f.Flags {
		if v.Value == nil {
			continue
		}

		prefix := "-"

		if v.Long {
			prefix += "-"
		}

		switch tv := v.Value.(type) {
		case bool:
			if explicitBool {
				script = fmt.Sprintf("%s %s%s%s%v", script, prefix, k, optionDelim, tv)
			} else {
				script = fmt.Sprintf("%s %s%s", script, prefix, k)
			}
		default:
			script = fmt.Sprintf("%s %s%s%s%v", script, prefix, k, optionDelim, tv)
		}
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
	return parseScript((f.Subcommands)[f.Choice], script, optionDelim, explicitBool)
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

	flags := map[string]*OptionValue{}
	args := map[string]*OptionValue{}
	subcommands := map[string]*Form{}

	f.Flags = flags
	f.Args = args
	f.Subcommands = subcommands

	for _, f := range c.Flags {
		flags[f.Name] = &OptionValue{
			Value: nil,
			Long:  f.Long,
		}
	}

	for _, a := range c.Args {
		args[a.Name] = &OptionValue{
			Value: nil,
		}
	}

	for _, c := range c.Subcommands {
		form := &Form{}
		parseForm(&c, form)
		f.Subcommands[c.Name] = form
	}
}
