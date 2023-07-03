package config

import "fmt"

type Form struct {
	Flags       map[string]any
	Args        map[string]any
	Subcommands map[string]*Form
	Choice      string
}

func (c CLI) Script(f Form) string {
	return parseScript(&f, c.Name, c.FlagDelim, c.DoubleDashesForLongFlags)
}

func parseScript(f *Form, script string, flagDelim string, doubleDashes bool) string {
	for k, v := range f.Flags {
		if v == nil {
			continue
		}

		prefix := "-"

		// single dash sometimes starts flags with longer names
		// Ref: https://pkg.go.dev/flag#hdr-Command_line_flag_syntax
		if doubleDashes && len(k) > 1 {
			prefix += "-"
		}

		script = fmt.Sprintf("%s %s%s%s%s", script, prefix, k, flagDelim, v)
	}

	for _, v := range f.Args {
		if v == nil {
			continue
		}
		script = fmt.Sprintf("%s %s", script, v)
	}

	if len(f.Subcommands) == 0 || f.Choice == "" {
		return script
	}

	script = fmt.Sprintf("%s %s", script, f.Choice)
	return parseScript((f.Subcommands)[f.Choice], script, flagDelim, doubleDashes)
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

	flags := map[string]any{}
	args := map[string]any{}
	subcommands := map[string]*Form{}

	f.Flags = flags
	f.Args = args
	f.Subcommands = subcommands

	for _, f := range c.Flags {
		flags[f.Name] = nil
	}

	for _, a := range c.Args {
		args[a.Name] = nil
	}

	for _, c := range c.Subcommands {
		form := &Form{}
		parseForm(&c, form)
		f.Subcommands[c.Name] = form
	}
}
