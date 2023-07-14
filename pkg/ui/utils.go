package ui

import (
	"CLI2UI/pkg/config"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
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

func updateCheckedOptions(v *map[string]*config.OptionValue, checked []string) {
	for k, v := range *v {
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
			v.ResetValue()
		}
	}
}

func updateValueEvent(field string, p Path, o config.Option) []sunmao.EventHandler {
	return []sunmao.EventHandler{
		{
			Type:        "onChange",
			ComponentId: "$utils",
			Method: sunmao.EventMethod{
				Name: "binding/v1/UpdateOptionValue",
				Parameters: UpdateOptionValueParams{
					OptionName: o.Name,
					Path:       p,
					Value:      fmt.Sprintf("{{ %s.%s }}", p.optionValueInputId(o.Name), field),
				},
			},
		},
	}
}

type Validator[T any] interface {
	Validate(T) bool
}

func (u UI) stringComponent(o config.Option, p Path) (sunmao.BaseComponentBuilder, []Validator[string]) {
	var comp sunmao.BaseComponentBuilder
	vs := []Validator[string]{}

	comp = u.arco.NewInput().
		Id(p.optionValueInputId(o.Name)).
		Properties(structToMap(InputProperties[string]{
			Size:     "default",
			Disabled: "{{ exec.state.isRunning }}",
		})).
		Event(updateValueEvent("value", p, o))

	switch o.Annotations.Format {
	case config.FormatAnnotationDate:
		comp = u.arco.NewDatePicker().
			Id(p.optionValueInputId(o.Name)).
			Properties(structToMap(DatePickerProperties[string]{
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(updateValueEvent("dateString", p, o))

	}

	// TODO(xinxi.guo): implement validation using annotations

	return comp, vs
}
