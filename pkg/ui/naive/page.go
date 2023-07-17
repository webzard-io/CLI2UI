package naive

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
	"fmt"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

func (u UI) buildPage() {
	cs := []sunmao.BaseComponentBuilder{
		u.layout(),
		u.helpModal(),
		u.dryRunModal(),
	}

	for _, c := range cs {
		u.arco.Component(c)
	}
}

func (u UI) layout() sunmao.BaseComponentBuilder {
	p := Path{}

	return u.arco.NewLayout().
		Properties(ui.StructToMap(ui.LayoutProperties{
			ShowHeader: true,
		})).
		Style("header", `
		display: flex;
		flex-direction: row;
		place-items: center;
		padding: 0.75rem 1rem;
		gap: 0.5rem;
		`).
		Style("content", `
		display: flex;
		flex-direction: column;
		place-items: center;
		padding: 0.75rem 1rem;
		gap: 0.5rem;
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"header": u.headerElements(),
			"content": {
				u.commandStack(p, u.cli.Command),
				u.runButton(),
				u.dryRunButton(),
				u.terminal(),
				u.stopButton(),
			},
		})
}

func (u UI) stopButton() sunmao.BaseComponentBuilder {
	return u.arco.NewButton().
		Properties(ui.StructToMap(ui.ButtonProperties[string]{
			Type:     "primary",
			Status:   "default",
			Size:     "default",
			Shape:    "square",
			Text:     "Stop",
			Disabled: "{{ !exec.state.isRunning }}",
		})).
		Style("content", "width: 100%;").
		Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "$utils",
				Method: sunmao.EventMethod{
					Name: "binding/v1/Stop",
				},
			},
		})
}

func (u UI) terminal() sunmao.BaseComponentBuilder {
	return u.arco.NewCollapse().
		Properties(ui.StructToMap(ui.CollapseProperties{
			DefaultActiveKey: []string{"0"},
			Options: []ui.CollapseItemProperties{
				{
					Key:            "0",
					Header:         "Result (Standard I/O)",
					ShowExpandIcon: true,
				},
			},
			ExpandIconPosition: "left",
			LazyLoad:           true,
		})).
		Style("content", `
		width: 100%;
		.arco-collapse-item-content-box { 
			background: white;
		}
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				u.arco.NewText().Content("Standard Output"),
				u.c2u.NewTerminal().Text("{{ exec.state.stdout }}"),
				u.arco.NewText().Content("Standard Error"),
				u.c2u.NewTerminal().Text("{{ exec.state.stderr }}"),
			},
		})
}

func (u UI) dryRunButton() sunmao.BaseComponentBuilder {
	return u.arco.NewButton().
		Properties(ui.StructToMap(ui.ButtonProperties[string]{
			Type:     "secondary",
			Status:   "default",
			Size:     "default",
			Shape:    "square",
			Text:     "Dry run",
			Disabled: "{{ exec.state.isRunning }}",
		})).
		Style("content", "width: 100%;").
		Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "$utils",
				Method: sunmao.EventMethod{
					Name: "binding/v1/DryRun",
				},
			},
			{
				Type:        "onClick",
				ComponentId: "DryRunModal",
				Method: sunmao.EventMethod{
					Name: "openModal",
				},
			},
		})
}

func (u UI) dryRunModal() sunmao.BaseComponentBuilder {
	copy := u.arco.NewButton().
		Properties(ui.StructToMap(ui.ButtonProperties[string]{
			Type:   "secondary",
			Status: "default",
			Size:   "default",
			Shape:  "square",
			Text:   "Copy",
		})).
		Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "$utils",
				Method: sunmao.EventMethod{
					Parameters: "{{ navigator.clipboard.writeText(dryRun.state); }}",
				},
			},
		})

	code := u.arco.NewText().Content("{{ `$ ${dryRun.state}` }}")

	close := u.arco.NewButton().
		Properties(ui.StructToMap(ui.ButtonProperties[string]{
			Type:   "default",
			Status: "default",
			Size:   "default",
			Shape:  "square",
			Text:   "Close",
		})).
		Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "DryRunModal",
				Method: sunmao.EventMethod{
					Name: "closeModal",
				},
			},
		})

	modal := u.arco.NewModal().Id("DryRunModal").
		Properties(ui.StructToMap(ui.ModalProperties{
			Title:         "Dry Run Result",
			Mask:          true,
			Closable:      true,
			MaskClosable:  true,
			UnmountOnExit: true,
		})).
		Style("content", `
		width: 80vw;
		.arco-modal-content {
			background: black;
			color: white;
			font-size: 1rem;
			font-family: monospace;
		}
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {code},
			"footer":  {copy, close},
		})

	return modal
}

func (u UI) runButton() sunmao.BaseComponentBuilder {
	return u.arco.NewButton().
		Properties(ui.StructToMap(ui.ButtonProperties[string]{
			Text:     "Run",
			Type:     "primary",
			Status:   "default",
			Size:     "default",
			Shape:    "square",
			Disabled: "{{ exec.state.isRunning }}",
		})).
		Style("content", "width: 100%;").
		Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "$utils",
				Method: sunmao.EventMethod{
					Name: "binding/v1/Run",
				},
			},
		})
}

func (u UI) commandStack(p Path, c config.Command) *sunmao.StackComponentBuilder {
	cs := []sunmao.BaseComponentBuilder{}
	if c.Description != "" {
		cs = append(cs, u.arco.NewText().
			Content(c.Description).
			Style("content", "color: var(--color-text-2)"))
	}
	cs = append(cs, u.optionsInputForm(p, c))

	if len(c.Subcommands) > 0 {
		cs = append(cs, u.subcommandsTab(p, c))
		for i, c := range c.Subcommands {
			pre := p.Append(c.Name)
			s := u.commandStack(Path{pre}, c).Slot(sunmao.Container{
				ID:   p.commandStackId(),
				Slot: "content",
			}, fmt.Sprintf("{{ %s.activeTab === %d }}", p.subcommandTabsId(), i))
			cs = append(cs, s)
		}
	}

	return u.arco.NewStack().
		Id(p.commandStackId()).
		Properties(ui.StructToMap(ui.StackProperties{
			Direction: "vertical",
			Spacing:   6,
		})).
		Style("content", "width: 100%;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": cs,
		})
}

func (u UI) subcommandsTab(p Path, c config.Command) sunmao.BaseComponentBuilder {
	tabs := []ui.TabProperties{}
	values := []string{}

	for _, c := range c.Subcommands {
		tabs = append(tabs, ui.TabProperties{
			Title:         c.DisplayName(),
			DestroyOnHide: true,
		})
		values = append(values, c.Name)
	}

	form := p.TraverseForm(u.fTpl)
	form.Choice = values[0]

	activeTab := fmt.Sprintf("{{ %s.activeTab }}", p.subcommandTabsId())

	return u.arco.NewTabs().
		Id(p.subcommandTabsId()).
		Properties(ui.StructToMap(ui.TabsProperties{
			Type:        "line",
			TabPosition: "top",
			Size:        "default",
			Tabs:        tabs,
		})).
		Style("content", "width: 100%;").
		Event([]sunmao.EventHandler{
			{
				Type:        "onChange",
				ComponentId: "$utils",
				Method: sunmao.EventMethod{
					Name: "binding/v1/UpdateSubcommand",
					Parameters: UpdateSubcommandParams[string]{
						SubcommandIndex: activeTab,
						Path:            p,
						Values:          values,
					},
				},
			},
		})
}

func (u UI) optionsInputForm(p Path, c config.Command) sunmao.BaseComponentBuilder {
	os, required, inputs := u.parseOptions(p, c)

	cb := u.arco.NewCheckbox().
		Id(p.optionsCheckboxId()).
		Properties(ui.StructToMap(ui.CheckboxProperties[string]{
			Options:              os,
			DefaultCheckedValues: required,
			Direction:            "horizontal",
			Disabled:             "{{ exec.state.isRunning }}",
		})).
		Style("content", `
		display: flex;
		align-items: flex-end;
		`).
		Event([]sunmao.EventHandler{
			{
				Type:        "onChange",
				ComponentId: "$utils",
				Method: sunmao.EventMethod{
					Name: "binding/v1/UpdateCheckedOptions",
					Parameters: UpdateCheckedOptionsParams[string]{
						Path:          p,
						CheckedValues: fmt.Sprintf("{{ %s.checkedValues }}", p.optionsCheckboxId()),
					},
				},
			},
		})

	cbWrapper := u.arco.NewStack().
		Properties(ui.StructToMap(ui.StackProperties{
			Direction: "horizontal",
			Justify:   "flex-end",
		})).Children(map[string][]sunmao.BaseComponentBuilder{
		"content": {cb},
	})

	contentElements := []sunmao.BaseComponentBuilder{}
	if len(os) > 0 {
		contentElements = append(contentElements, cbWrapper)
	}
	contentElements = append(contentElements, inputs...)

	s := u.arco.NewStack().
		Id(p.optionValuesFormId()).
		Properties(ui.StructToMap(ui.StackProperties{
			Direction: "vertical",
			Spacing:   6,
		})).
		Style("content", "width: 100%;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": contentElements,
		})

	return s
}

func (u UI) parseOptions(p Path, c config.Command) ([]ui.CheckboxOptionProperties, []string, []sunmao.BaseComponentBuilder) {
	os := []ui.CheckboxOptionProperties{}
	required := []string{}
	inputs := []sunmao.BaseComponentBuilder{}

	for _, f := range c.Flags {
		inputs = append(inputs, u.optionInput(p, f))
		os = append(os, ui.CheckboxOptionProperties{
			Label:    f.DisplayName(),
			Value:    f.Name,
			Disabled: f.Required,
		})

		if f.Required {
			required = append(required, f.Name)
		}
	}

	for _, a := range c.Args {
		inputs = append(inputs, u.optionInput(p, a))
		os = append(os, ui.CheckboxOptionProperties{
			Label:    a.DisplayName(),
			Value:    a.Name,
			Disabled: a.Required,
		})

		if a.Required {
			required = append(required, a.Name)
		}
	}

	return os, required, inputs
}

func (u UI) optionInput(p Path, o config.Option) sunmao.BaseComponentBuilder {
	cs := []sunmao.BaseComponentBuilder{u.inputType(p, o)}
	if o.Description != "" {
		cs = append(cs, u.arco.NewText().
			Content(o.Description).
			Style("content", "color: var(--color-text-2);"))
	}

	return u.arco.NewFormControl().
		Properties(ui.StructToMap(ui.FormControlProperties{
			Label: ui.TextProperties{
				Format: "plain",
				Raw:    o.DisplayName(),
			},
			Layout:     "horizontal",
			Required:   o.Required,
			LabelAlign: "left",
			LabelCol: ui.ColumnProperties{
				Span: 4,
			},
			WrapperCol: ui.ColumnProperties{
				Span: 20,
			},
		})).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": cs,
		}).
		Slot(sunmao.Container{
			ID:   p.optionValuesFormId(),
			Slot: "content",
		}, fmt.Sprintf("{{ %s.checkedValues.some(o => o === \"%s\") }}", p.optionsCheckboxId(), o.Name))
}

func (u UI) inputType(p Path, o config.Option) sunmao.BaseComponentBuilder {
	// TODO(xinxi.guo): typeComponent() will ultimately replace all these
	switch o.Type {
	case config.OptionTypeNumber:
		return u.arco.NewNumberInput().
			Id(p.optionValueInputId(o.Name)).
			Properties(ui.StructToMap(ui.NumberInputProperties[string]{
				Size:     "default",
				Max:      99,
				Step:     1,
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(updateValueEvent("value", p, o))
	case config.OptionTypeArray:
		return u.c2u.NewArrayInput().
			Id(p.optionValueInputId(o.Name)).
			Properties(ui.StructToMap(ui.ArrayInputProperties[string]{
				Value:    []string{""},
				Type:     "string",
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(updateValueEvent("value", p, o))
	case config.OptionTypeBoolean:
		return u.arco.NewSwitch().
			Id(p.optionValueInputId(o.Name)).
			Properties(ui.StructToMap(ui.SwitchProperties[string]{
				Type:     "circle",
				Size:     "default",
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(updateValueEvent("value", p, o))
	case config.OptionTypeEnum:
		options := []ui.SelectOptionProperties{}
		for _, o := range o.Options {
			options = append(options, ui.SelectOptionProperties{
				Text:  o,
				Value: o,
			})
		}
		return u.arco.NewSelect().
			Id(p.optionValueInputId(o.Name)).
			Properties(ui.StructToMap(ui.SelectProperties[string]{
				Bordered:            true,
				UnmountOnExit:       true,
				Options:             options,
				Size:                "default",
				AutoAlignPopupWidth: true,
				Position:            "bottom",
				MountToBody:         true,
				Disabled:            "{{ exec.state.isRunning }}",
			})).
			Event(updateValueEvent("value", p, o))
	}

	// TODO(xinxi.guo): implement validation
	comp, _ := u.stringComponent(o, p)
	return comp
}

func (u UI) headerElements() []sunmao.BaseComponentBuilder {
	title := u.arco.NewText().
		Content(u.cli.Name).Style("content",
		`
		font-size: 1.25rem;
		font-weight: bold;
		`)

	help := u.arco.NewButton().
		Properties(ui.StructToMap(ui.ButtonProperties[string]{
			Shape: "square",
			Text:  "Help",
			Type:  "default",
		})).
		Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "HelpModal",
				Method: sunmao.EventMethod{
					Name: "openModal",
				},
			},
		})

	return []sunmao.BaseComponentBuilder{title, help}
}

func (u UI) helpModal() sunmao.BaseComponentBuilder {
	help := u.c2u.NewTextDisplay().
		Style("content", `
		height: 24rem;
		overflow: scroll;
		`).
		Content(ui.TextDisplayProperties{
			Text:   u.cli.Help,
			Format: "code",
		})

	close := u.arco.NewButton().
		Properties(ui.StructToMap(ui.ButtonProperties[string]{
			Type:   "default",
			Status: "default",
			Size:   "default",
			Shape:  "square",
			Text:   "Close",
		})).
		Event([]sunmao.EventHandler{
			{
				Type:        "onClick",
				ComponentId: "HelpModal",
				Method: sunmao.EventMethod{
					Name: "closeModal",
				},
			},
		})

	modal := u.arco.NewModal().Id("HelpModal").
		Properties(ui.StructToMap(ui.ModalProperties{
			Title:         "Help",
			Mask:          true,
			Closable:      true,
			MaskClosable:  true,
			UnmountOnExit: true,
		})).
		Style("content", "width: 80vw;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {help},
			"footer":  {close},
		})

	return modal
}