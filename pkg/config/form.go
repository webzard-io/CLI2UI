package config

import (
	"fmt"
	"sort"
)

type Form struct {
	Flags       map[string]*OptionValue
	Args        map[string]*OptionValue
	Subcommands map[string]*Form
	Choice      string
}

type OptionValue struct {
	Value        any
	Long         bool
	Enabled      bool
	required     bool
	defaultValue any
}

func (c CLI) Script(f Form) string {
	return parseScript(&f, c.Command.Name, c.OptionDelim, c.ExplicitBool)
}

func parseScript(f *Form, script string, optionDelim string, explicitBool bool) string {
	for _, k := range orderedKeys(f.Flags) {
		v := f.Flags[k]
		if !v.Enabled {
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

	for _, k := range orderedKeys(f.Args) {
		v := f.Args[k]
		if !v.Enabled {
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

func orderedKeys(m map[string]*OptionValue) []string {
	keys := []string{}

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
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
		dv := f.Default
		if dv == nil {
			// TODO(xinxi.guo): type system has to be enhanced to make use of `Option.Default`, this is a workaround for now
			dv = fmt.Sprintf("<%s>", f.Name)
		}
		flags[f.Name] = &OptionValue{
			Long:         f.Long,
			Value:        dv,
			required:     f.Required,
			Enabled:      f.Required,
			defaultValue: dv,
		}
	}

	for _, a := range c.Args {
		dv := a.Default
		if dv == nil {
			dv = fmt.Sprintf("<%s>", a.Name)
		}
		args[a.Name] = &OptionValue{
			Value:        fmt.Sprintf("<%s>", a.Name),
			required:     a.Required,
			Enabled:      a.Required,
			defaultValue: dv,
		}
	}

	for _, c := range c.Subcommands {
		form := &Form{}
		parseForm(&c, form)
		f.Subcommands[c.Name] = form
	}
}

func (f Form) Clear() {
	for _, v := range f.Args {
		v.Value = v.defaultValue
		v.Enabled = v.required
	}

	for _, v := range f.Flags {
		v.Value = v.defaultValue
		v.Enabled = v.required
	}

	for _, k := range f.Subcommands {
		k.Clear()
	}
}

func (f Form) Clone() *Form {
	t := &Form{
		Flags:       map[string]*OptionValue{},
		Args:        map[string]*OptionValue{},
		Subcommands: map[string]*Form{},
		Choice:      f.Choice,
	}

	for k, v := range f.Flags {
		t.Flags[k] = &OptionValue{
			Value:        v.Value,
			Long:         v.Long,
			Enabled:      v.Enabled,
			required:     v.required,
			defaultValue: v.defaultValue,
		}
	}

	for k, v := range f.Args {
		t.Args[k] = &OptionValue{
			Value:        v.Value,
			Long:         v.Long,
			Enabled:      v.Enabled,
			required:     v.required,
			defaultValue: v.defaultValue,
		}
	}

	for k, v := range f.Subcommands {
		t.Subcommands[k] = v.Clone()
	}

	return t
}
