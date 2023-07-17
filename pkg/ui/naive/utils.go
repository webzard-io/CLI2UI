package naive

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
	"fmt"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type Path struct {
	ui.Path
}

func (p Path) optionsCheckboxId() string {
	return fmt.Sprintf("%sOptionsCheckbox", p)
}

func (p Path) optionValuesFormId() string {
	return fmt.Sprintf("%sOptionValuesForm", p)
}

func (p Path) optionValueInputId(option string) string {
	return fmt.Sprintf("%s%sOptionValueInput", p, ui.KebabToPascalCase(option))
}

func (p Path) subcommandTabsId() string {
	return fmt.Sprintf("%sSubcommandTabs", p)
}

func (p Path) commandStackId() string {
	return fmt.Sprintf("%sCommandStack", p)
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
		Properties(ui.StructToMap(ui.InputProperties[string]{
			Size:     "default",
			Disabled: "{{ exec.state.isRunning }}",
		})).
		Event(updateValueEvent("value", p, o))

	switch o.Annotations.Format {
	case config.FormatAnnotationDate:
		comp = u.arco.NewDatePicker().
			Id(p.optionValueInputId(o.Name)).
			Properties(ui.StructToMap(ui.DatePickerProperties[string]{
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(updateValueEvent("dateString", p, o))

	}

	// TODO(xinxi.guo): implement validation using annotations

	return comp, vs
}
