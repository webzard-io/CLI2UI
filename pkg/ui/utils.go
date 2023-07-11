package ui

import (
	"CLI2UI/pkg/config"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func structToMap(s interface{}) map[string]interface{} {
	data, err := json.Marshal(s)
	if err != nil {
		log.Fatal("error marshaling struct to JSON", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal("error unmarshaling JSON to map", err)
	}

	return result
}

func kebabToPascalCase(i string) string {
	words := strings.Split(i, "-")
	for i, word := range words {
		words[i] = cases.Title(language.English).String(word)
	}
	return strings.Join(words, "")
}

type Path []string

func (p Path) String() string {
	c := []string{}
	for _, s := range p {
		c = append(c, kebabToPascalCase(s))
	}

	return strings.Join(c, "")
}

func (p Path) append(s string) Path {
	return append(p, s)
}

func (p Path) optionsCheckboxId() string {
	return fmt.Sprintf("%sOptionsCheckbox", p)
}

func (p Path) optionValuesFormId() string {
	return fmt.Sprintf("%sOptionValuesForm", p)
}

func (p Path) optionValueInputId(option string) string {
	return fmt.Sprintf("%s%sOptionValueInput", p, kebabToPascalCase(option))
}

func (p Path) subcommandTabsId() string {
	return fmt.Sprintf("%sSubcommandTabs", p)
}

func (p Path) commandStackId() string {
	return fmt.Sprintf("%sCommandStack", p)
}

func toStruct[T any](s any) T {
	var t T

	b, _ := json.Marshal(s)
	_ = json.Unmarshal(b, &t)

	return t
}

func (p Path) traverseForm(f *config.Form) *config.Form {
	form := f
	for _, c := range p {
		form = form.Subcommands[c]
	}
	return form
}

func clearForm(f *config.Form) {
	for k := range f.Args {
		f.Args[k].Value = fmt.Sprintf("<%s>", k)
		f.Args[k].Enabled = f.Args[k].Required
	}

	for k := range f.Flags {
		f.Flags[k].Value = fmt.Sprintf("<%s>", k)
		f.Flags[k].Enabled = f.Flags[k].Required
	}

	for k := range f.Subcommands {
		clearForm(f.Subcommands[k])
	}
}

func updateCheckedOptions(f *config.Form, checked []string) {
	for k, v := range f.Flags {
		found := false
		for _, cv := range checked {
			if k == cv {
				v.Enabled = true
				found = true
				break
			}
		}
		if !found {
			v.Enabled = false
		}
	}

	for k, v := range f.Args {
		found := false
		for _, cv := range checked {
			if k == cv {
				v.Enabled = true
				found = true
				break
			}
		}
		if !found {
			v.Enabled = false
		}
	}
}
