package ui

import (
	"CLI2UI/pkg/config"
	"fmt"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

func (u UI) buildPage() {
	cs := []sunmao.BaseComponentBuilder{
		u.layout(),
		u.helpModal(),
	}

	for _, c := range cs {
		u.arco.Component(c)
	}
}

func (u UI) layout() sunmao.BaseComponentBuilder {
	p := Path{}

	return u.arco.NewLayout().
		Properties(structToMap(LayoutProperties{
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
		Properties(structToMap(ButtonProperties[string]{
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
		Properties(structToMap(CollapseProperties{
			DefaultActiveKey: []string{"0"},
			Options: []CollapseItemProperties{
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
		Properties(structToMap(ButtonProperties[string]{
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
		})
}

func (u UI) runButton() sunmao.BaseComponentBuilder {
	return u.arco.NewButton().
		Properties(structToMap(ButtonProperties[string]{
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
	cs := []sunmao.BaseComponentBuilder{
		u.arco.NewText().
			Content(fmt.Sprintf("Command description: %s", c.Description)).
			Style("content", "color: var(--color-text-2)"),
		u.optionsInputForm(p, c),
	}

	if len(c.Subcommands) > 0 {
		cs = append(cs, u.subcommandsTab(p, c))
		for i, c := range c.Subcommands {
			pre := p.append(c.Name)
			s := u.commandStack(pre, c).Slot(sunmao.Container{
				ID:   p.commandStackId(),
				Slot: "content",
			}, fmt.Sprintf("{{ %s.activeTab === %d }}", p.subcommandTabsId(), i))
			cs = append(cs, s)
		}
	}

	return u.arco.NewStack().
		Id(p.commandStackId()).
		Properties(structToMap(StackProperties{
			Direction: "vertical",
			Spacing:   6,
		})).
		Style("content", "width: 100%;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": cs,
		})
}

func (u UI) subcommandsTab(p Path, c config.Command) sunmao.BaseComponentBuilder {
	tabs := []TabProperties{}

	// TODO(xinxi.guo): use `c.DisplayName()` instead once possible
	for _, c := range c.Subcommands {
		tabs = append(tabs, TabProperties{
			Title:         c.Name,
			DestroyOnHide: true,
		})
	}

	defaultTab := tabs[0].Title
	form := p.traverseForm(u.fTpl)
	form.Choice = defaultTab

	activeTab := fmt.Sprintf("{{ %s.activeTab }}", p.subcommandTabsId())

	return u.arco.NewTabs().
		Id(p.subcommandTabsId()).
		Properties(structToMap(TabsProperties{
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
						Tabs:            tabs,
					},
				},
			},
		})
}

func (u UI) optionsInputForm(p Path, c config.Command) sunmao.BaseComponentBuilder {
	os, required, inputs := u.parseOptions(p, c)

	// TODO(xinxi.guo): notify the user when there is no option available
	cb := u.arco.NewCheckbox().
		Id(p.optionsCheckboxId()).
		Properties(structToMap(CheckboxProperties[string]{
			Options:              os,
			DefaultCheckedValues: required,
			Direction:            "horizontal",
			Disabled:             "{{ exec.state.isRunning }}",
		})).
		Style("content", `
		display: flex;
		align-items: flex-end;
		`)

	cbWrapper := u.arco.NewStack().
		Properties(structToMap(StackProperties{
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
		Properties(structToMap(StackProperties{
			Direction: "vertical",
			Spacing:   6,
		})).
		Style("content", "width: 100%;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": contentElements,
		})

	return s
}

func (u UI) parseOptions(p Path, c config.Command) ([]CheckboxOptionProperties, []string, []sunmao.BaseComponentBuilder) {
	os := []CheckboxOptionProperties{}
	required := []string{}
	inputs := []sunmao.BaseComponentBuilder{}

	for _, f := range c.Flags {
		inputs = append(inputs, u.optionInput(p, f))
		os = append(os, CheckboxOptionProperties{
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
		os = append(os, CheckboxOptionProperties{
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
	return u.arco.NewFormControl().
		Properties(structToMap(FormControlProperties{
			Label: TextProperties{
				Format: "plain",
				Raw:    o.DisplayName(),
			},
			Layout:     "horizontal",
			Required:   o.Required,
			LabelAlign: "left",
			LabelCol: ColumnProperties{
				Span: 4,
			},
			WrapperCol: ColumnProperties{
				Span: 20,
			},
		})).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				u.inputType(p, o),
				u.arco.NewText().
					Content(fmt.Sprintf("Option description: %s", o.Description)).
					Style("content", "color: var(--color-text-2);"),
			},
		}).
		Slot(sunmao.Container{
			ID:   p.optionValuesFormId(),
			Slot: "content",
		}, fmt.Sprintf("{{ %s.checkedValues.some(o => o === \"%s\") }}", p.optionsCheckboxId(), o.Name))
}

func (u UI) inputType(p Path, o config.Option) sunmao.BaseComponentBuilder {
	es := []sunmao.EventHandler{
		{
			Type:        "onChange",
			ComponentId: "$utils",
			Method: sunmao.EventMethod{
				Name: "binding/v1/UpdateOptionValue",
				Parameters: UpdateOptionValueParams{
					OptionName: o.Name,
					Path:       p,
					Value:      fmt.Sprintf("{{ %s.value }}", p.optionValueInputId(o.Name)),
				},
			},
		},
	}

	switch o.Type {
	case config.OptionTypeNumber:
		return u.arco.NewNumberInput().
			Id(p.optionValueInputId(o.Name)).
			Properties(structToMap(NumberInputProperties[string]{
				Size:     "default",
				Max:      99,
				Step:     1,
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(es)
	case config.OptionTypeArray:
		return u.c2u.NewArrayInput().
			Id(p.optionValueInputId(o.Name)).
			Properties(structToMap(ArrayInputProperties[string]{
				Value:    []string{""},
				Type:     "string",
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(es)
	case config.OptionTypeBoolean:
		return u.arco.NewSwitch().
			Id(p.optionValueInputId(o.Name)).
			Properties(structToMap(SwitchProperties[string]{
				Type:     "circle",
				Size:     "default",
				Disabled: "{{ exec.state.isRunning }}",
			})).
			Event(es)
	case config.OptionTypeEnum:
		options := []SelectOptionProperties{}
		for _, o := range o.Options {
			options = append(options, SelectOptionProperties{
				Text:  o,
				Value: o,
			})
		}
		return u.arco.NewSelect().
			Id(p.optionValueInputId(o.Name)).
			Properties(structToMap(SelectProperties[string]{
				Bordered:            true,
				UnmountOnExit:       true,
				Options:             options,
				Size:                "default",
				AutoAlignPopupWidth: true,
				Position:            "bottom",
				MountToBody:         true,
				Disabled:            "{{ exec.state.isRunning }}",
			})).
			Event(es)
	}

	// TODO(xinxi.guo): implement validation
	return u.arco.NewInput().
		Id(p.optionValueInputId(o.Name)).
		Properties(structToMap(InputProperties[string]{
			Size:     "default",
			Disabled: "{{ exec.state.isRunning }}",
		})).
		Event(es)
}

func (u UI) headerElements() []sunmao.BaseComponentBuilder {
	title := u.arco.NewText().
		Content(u.cli.Name).Style("content",
		`
		font-size: 1.25rem;
		font-weight: bold;
		`)

	help := u.arco.NewButton().
		Properties(structToMap(ButtonProperties[string]{
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
		Content(TextDisplayProperties{
			Text:   u.cli.Help,
			Format: "code",
		})

	close := u.arco.NewButton().
		Properties(structToMap(ButtonProperties[string]{
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
		Properties(structToMap(ModalProperties{
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
