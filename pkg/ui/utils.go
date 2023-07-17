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

func StructToMap(s interface{}) map[string]interface{} {
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

func KebabToPascalCase(i string) string {
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
		c = append(c, KebabToPascalCase(s))
	}

	return strings.Join(c, "")
}

func (p Path) Append(s string) Path {
	return append(p, s)
}

func (p Path) OptionsCheckboxId() string {
	return fmt.Sprintf("%sOptionsCheckbox", p)
}

func (p Path) OptionValuesFormId() string {
	return fmt.Sprintf("%sOptionValuesForm", p)
}

func (p Path) OptionValueInputId(option string) string {
	return fmt.Sprintf("%s%sOptionValueInput", p, KebabToPascalCase(option))
}

func (p Path) CommandStackId() string {
	return fmt.Sprintf("%sCommandStack", p)
}

func ToStruct[T any](s any) T {
	var t T

	b, _ := json.Marshal(s)
	_ = json.Unmarshal(b, &t)

	return t
}

func (p Path) TraverseForm(f *config.Form) *config.Form {
	return p.TraverseFormWithCallback(f, nil)
}

func (p Path) TraverseFormWithCallback(f *config.Form, cb func(string, *config.Form)) *config.Form {
	form := f
	for _, c := range p {
		if cb != nil {
			cb(c, form)
		}
		form = form.Subcommands[c]
	}
	return form
}
