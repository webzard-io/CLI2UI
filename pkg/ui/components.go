package ui

import "github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"

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
	b.Properties(structToMap(v))
	return b
}

type CheckboxMenuBuilder struct {
	*sunmao.InnerComponentBuilder[*CheckboxMenuBuilder]
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

func (b *CLI2UIAppBuilder) NewArrayInput() *ArrayInputBuilder {
	t := &ArrayInputBuilder{
		InnerComponentBuilder: sunmao.NewInnerComponent[*ArrayInputBuilder](b.AppBuilder),
	}
	t.Inner = t
	return t.Type("cli2ui/v1/arrayInput")
}
