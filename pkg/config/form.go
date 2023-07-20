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
	format      Format
}

type OptionValue struct {
	Value        any
	Long         bool
	Enabled      bool
	required     bool
	defaultValue any
	index        int
	name         string
}

func (c CLI) Script(f Form) (string, Format) {
	return parseScript(&f, c.Command.Name, c.OptionDelim, c.ExplicitBool)
}

func parseScript(f *Form, script string, optionDelim string, explicitBool bool) (string, Format) {
	as, fs := sortedOptions(f)
	for _, v := range fs {
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
				script = fmt.Sprintf("%s %s%s%s%v", script, prefix, v.name, optionDelim, tv)
			} else {
				if !v.Value.(bool) {
					continue
				}
				script = fmt.Sprintf("%s %s%s", script, prefix, v.name)
			}
		default:
			script = fmt.Sprintf("%s %s%s%s%v", script, prefix, v.name, optionDelim, tv)
		}
	}

	for _, v := range as {
		if !v.Enabled {
			continue
		}
		script = fmt.Sprintf("%s %s", script, v.Value)
	}

	if len(f.Subcommands) == 0 || f.Choice == "" {
		return script, f.format
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
	f.format = c.Format

	for _, f := range c.Flags {
		dv := f.Default
		if dv == nil {
			// TODO(xinxi.guo): type system has to be enhanced to make use of `Option.Default`, this is a workaround for now
			if f.Type == OptionTypeBoolean {
				dv = false
			} else {
				dv = fmt.Sprintf("<%s>", f.DisplayName())
			}
		}
		flags[f.Name] = &OptionValue{
			Long:         f.Long,
			Value:        dv,
			required:     f.Required,
			Enabled:      f.Required,
			defaultValue: dv,
			name:         f.Name,
		}
	}

	for i, a := range c.Args {
		dv := a.Default
		if dv == nil {
			// TODO(xinxi.guo): type system has to be enhanced to make use of `Option.Default`, this is a workaround for now
			if a.Type == OptionTypeBoolean {
				dv = false
			} else {
				dv = fmt.Sprintf("<%s>", a.DisplayName())
			}
		}
		args[a.Name] = &OptionValue{
			Value:        dv,
			required:     a.Required,
			Enabled:      a.Required,
			defaultValue: dv,
			name:         a.Name,
			index:        i,
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
		format:      f.format,
	}

	for k, v := range f.Flags {
		t.Flags[k] = &OptionValue{
			Value:        v.Value,
			Long:         v.Long,
			Enabled:      v.Enabled,
			required:     v.required,
			defaultValue: v.defaultValue,
			name:         v.name,
			index:        v.index,
		}
	}

	for k, v := range f.Args {
		t.Args[k] = &OptionValue{
			Value:        v.Value,
			Long:         v.Long,
			Enabled:      v.Enabled,
			required:     v.required,
			defaultValue: v.defaultValue,
			name:         v.name,
			index:        v.index,
		}
	}

	for k, v := range f.Subcommands {
		t.Subcommands[k] = v.Clone()
	}

	return t
}

func (o *OptionValue) ResetValue() {
	o.Value = o.defaultValue
}

type optionValueSorter []*OptionValue

func (o optionValueSorter) Len() int {
	return len(o)
}

func (o optionValueSorter) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o optionValueSorter) Less(i, j int) bool {
	return o[i].index >= o[j].index && o[i].name < o[j].name
}

func sortedOptions(f *Form) ([]*OptionValue, []*OptionValue) {
	as := []*OptionValue{}
	fs := []*OptionValue{}

	for _, a := range f.Args {
		as = append(as, a)
	}

	for _, f := range f.Flags {
		fs = append(fs, f)
	}

	sort.Sort(optionValueSorter(as))
	sort.Sort(optionValueSorter(fs))

	return as, fs
}
