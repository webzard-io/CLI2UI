package ui

import (
	"CLI2UI/pkg/config"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type CLI2UIAppBuilder struct {
	*sunmao.AppBuilder
}

func NewCLI2UIApp(appBuilder *sunmao.AppBuilder) *CLI2UIAppBuilder {
	b := &CLI2UIAppBuilder{
		appBuilder,
	}
	return b
}

type TextDisplayComponentBuilder struct {
	*sunmao.InnerComponentBuilder[*TextDisplayComponentBuilder]
}

type TextDisplayProperties struct {
	Text   string `json:"text"`
	Format string `json:"format"`
}

func (b *CLI2UIAppBuilder) NewTextDisplay() *TextDisplayComponentBuilder {
	t := &TextDisplayComponentBuilder{
		InnerComponentBuilder: sunmao.NewInnerComponent[*TextDisplayComponentBuilder](b.AppBuilder),
	}
	t.Inner = t
	return t.Type("cli2ui/v1/TextDisplay")
}

func (b *TextDisplayComponentBuilder) Content(v TextDisplayProperties) *TextDisplayComponentBuilder {
	b.Properties(StructToMap(v))
	return b
}

type CheckboxMenuBuilder struct {
	*sunmao.InnerComponentBuilder[*CheckboxMenuBuilder]
}

type CheckboxMenuOptionProperties struct {
	Label    string `json:"label"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
}

type CheckboxMenuProperties struct {
	Value   []string                       `json:"value"`
	Text    string                         `json:"text"`
	Options []CheckboxMenuOptionProperties `json:"options"`
}

func (b *CLI2UIAppBuilder) NewCheckboxMenu() *CheckboxMenuBuilder {
	t := &CheckboxMenuBuilder{
		InnerComponentBuilder: sunmao.NewInnerComponent[*CheckboxMenuBuilder](b.AppBuilder),
	}
	t.Inner = t
	return t.Type("cli2ui/v1/checkboxMenu")
}

type ResultBuilder struct {
	*sunmao.InnerComponentBuilder[*ResultBuilder]
}

func (b *CLI2UIAppBuilder) NewResult() *ResultBuilder {
	t := &ResultBuilder{
		InnerComponentBuilder: sunmao.NewInnerComponent[*ResultBuilder](b.AppBuilder),
	}
	t.Inner = t
	return t.Type("cli2ui/v1/result")
}

func (b *ResultBuilder) Data(value interface{}) *ResultBuilder {
	b.Properties(map[string]interface{}{
		"data": value,
	})
	return b
}

type TerminalBuilder struct {
	*sunmao.InnerComponentBuilder[*TerminalBuilder]
}

func (b *CLI2UIAppBuilder) NewTerminal() *TerminalBuilder {
	t := &TerminalBuilder{
		InnerComponentBuilder: sunmao.NewInnerComponent[*TerminalBuilder](b.AppBuilder),
	}
	t.Inner = t
	return t.Type("cli2ui/v1/terminal")
}

func (b *TerminalBuilder) Text(value interface{}) *TerminalBuilder {
	b.Properties(map[string]interface{}{
		"text": value,
	})
	return b
}

type ArrayInputBuilder struct {
	*sunmao.InnerComponentBuilder[*ArrayInputBuilder]
}

type ArrayInputProperties[T bool | string] struct {
	Value       []string `json:"value"`
	Type        string   `json:"type"`
	Placeholder string   `json:"placeholder"`
	Disabled    T        `json:"disabled"`
}

func (b *CLI2UIAppBuilder) NewArrayInput() *ArrayInputBuilder {
	t := &ArrayInputBuilder{
		InnerComponentBuilder: sunmao.NewInnerComponent[*ArrayInputBuilder](b.AppBuilder),
	}
	t.Inner = t
	return t.Type("cli2ui/v1/arrayInput")
}

type Validator[T any] interface {
	Validate(T) bool
}

func (u BaseUI) stringComponent(o config.Option, p Path) (sunmao.BaseComponentBuilder, []Validator[string]) {
	var comp sunmao.BaseComponentBuilder
	vs := []Validator[string]{}

	comp = u.Arco.NewInput().
		Id(p.OptionValueInputId(o.Name)).
		Properties(StructToMap(InputProperties[string]{
			Size:     "default",
			Disabled: "{{ exec.state.isRunning }}",
		})).
		Event(UpdateValueEvent("value", p, o))

	switch o.Annotations.Format {
	case config.FormatAnnotationDate:
		comp = u.Arco.NewDatePicker().
			Id(p.OptionValueInputId(o.Name)).
			Properties(StructToMap(DatePickerProperties[string]{
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(UpdateValueEvent("dateString", p, o))

	}

	// TODO(xinxi.guo): implement validation using annotations

	return comp, vs
}

func (u BaseUI) InputType(p Path, o config.Option) sunmao.BaseComponentBuilder {
	// TODO(xinxi.guo): typeComponent() will ultimately replace all these
	switch o.Type {
	case config.OptionTypeNumber:
		return u.Arco.NewNumberInput().
			Id(p.OptionValueInputId(o.Name)).
			Properties(StructToMap(NumberInputProperties[string]{
				Size:     "default",
				Max:      99,
				Step:     1,
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(UpdateValueEvent("value", p, o))
	case config.OptionTypeArray:
		return u.C2U.NewArrayInput().
			Id(p.OptionValueInputId(o.Name)).
			Properties(StructToMap(ArrayInputProperties[string]{
				Value:    []string{""},
				Type:     "string",
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(UpdateValueEvent("value", p, o))
	case config.OptionTypeBoolean:
		return u.Arco.NewSwitch().
			Id(p.OptionValueInputId(o.Name)).
			Properties(StructToMap(SwitchProperties[string]{
				Type:     "circle",
				Size:     "default",
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(UpdateValueEvent("value", p, o))
	case config.OptionTypeEnum:
		options := []SelectOptionProperties{}
		for _, o := range o.Options {
			options = append(options, SelectOptionProperties{
				Text:  o,
				Value: o,
			})
		}
		return u.Arco.NewSelect().
			Id(p.OptionValueInputId(o.Name)).
			Properties(StructToMap(SelectProperties[string]{
				Bordered:            true,
				UnmountOnExit:       true,
				Options:             options,
				Size:                "default",
				AutoAlignPopupWidth: true,
				Position:            "bottom",
				MountToBody:         true,
				Disabled:            "{{ exec.state.isRunning }}",
			})).
			Event(UpdateValueEvent("value", p, o))
	}

	// TODO(xinxi.guo): implement validation
	comp, _ := u.stringComponent(o, p)
	return comp
}
