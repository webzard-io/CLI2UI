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
				u.commandStack(u.cli.Command),
			},
		})
}

func (u UI) commandStack(c config.Command) sunmao.BaseComponentBuilder {
	cs := []sunmao.BaseComponentBuilder{u.optionsInputForm(c)}
	if len(c.Subcommands) > 0 {
		cs = append(cs, u.subcommandsTab(c))
	}
	return u.arco.NewStack().
		Properties(structToMap(StackProperties{
			Direction: "vertical",
			Spacing:   6,
		})).
		Style("content", "width: 100%;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": cs,
		})
}

func (u UI) subcommandsTab(c config.Command) sunmao.BaseComponentBuilder {
	tabs := []TabProperties{}

	for _, c := range c.Subcommands {
		tabs = append(tabs, TabProperties{
			Title:         c.Name,
			DestroyOnHide: true,
		})
	}

	return u.arco.NewTabs().
		Properties(structToMap(TabsProperties{
			Type:        "line",
			TabPosition: "top",
			Size:        "default",
			Tabs:        tabs,
		})).
		Style("content", "width: 100%;")
}

func (u UI) optionsInputForm(c config.Command) sunmao.BaseComponentBuilder {
	os, required, inputs := u.parseOptions(c)

	cb := u.c2u.NewCheckboxMenu().
		Id(optionsCheckboxId(c.Name)).
		Properties(structToMap(CheckboxMenuProperties{
			Value:   required,
			Text:    "Options",
			Options: os,
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

	contentElements := []sunmao.BaseComponentBuilder{cbWrapper}
	contentElements = append(contentElements, inputs...)

	s := u.arco.NewStack().
		Id(optionValuesFormId(c.Name)).
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

func (u UI) parseOptions(c config.Command) ([]CheckboxOption, []string, []sunmao.BaseComponentBuilder) {
	os := []CheckboxOption{}
	required := []string{}
	inputs := []sunmao.BaseComponentBuilder{}

	for _, f := range c.Flags {
		inputs = append(inputs, u.optionInput(c.Name, f))
		os = append(os, CheckboxOption{
			Label:    f.Name,
			Value:    f.Name,
			Disabled: f.Required,
		})

		if f.Required {
			required = append(required, f.Name)
		}
	}

	for _, a := range c.Args {
		inputs = append(inputs, u.optionInput(c.Name, a))
		os = append(os, CheckboxOption{
			Label:    a.Name,
			Value:    a.Name,
			Disabled: a.Required,
		})

		if a.Required {
			required = append(required, a.Name)
		}
	}

	return os, required, inputs
}

func (u UI) optionInput(cmd string, o config.FlagOrArg) sunmao.BaseComponentBuilder {
	return u.arco.NewFormControl().
		Id(optionValueInputId(cmd, o.Name)).
		Properties(structToMap(FormControlProperties{
			Label: TextProperties{
				Format: "plain",
				Raw:    o.Name,
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
			"content": {u.inputType(cmd, o)},
		}).
		Slot(sunmao.Container{
			ID:   optionValuesFormId(cmd),
			Slot: "content",
		}, fmt.Sprintf("{{ %s.value.some(o => o === \"%s\") }}", optionsCheckboxId(cmd), o.Name))
}

func (u UI) inputType(cmd string, o config.FlagOrArg) sunmao.BaseComponentBuilder {
	// TODO(xinxi.guo): disable when command is running
	switch o.Type {
	case config.FlagArgTypeNumber:
		return u.arco.NewNumberInput().
			Properties(structToMap(NumberInputProperties{
				DefaultValue: 1,
				Placeholder:  o.Default,
				Size:         "default",
				Max:          99,
				Step:         1,
			}))
	case config.FlagArgTypeArray:
		return u.c2u.NewArrayInput().
			Properties(structToMap(ArrayInputProperties{
				Value:       []string{""},
				Type:        "string",
				Placeholder: o.Default,
			}))
	case config.FlagArgTypeBoolean:
		return u.arco.NewSwitch().
			Properties(structToMap(SwitchProperties{
				Type: "circle",
				Size: "default",
			}))
	case config.FlagArgTypeEnum:
		options := []SelectOptionProperties{}
		for _, o := range o.Options {
			options = append(options, SelectOptionProperties{
				Text:  o,
				Value: o,
			})
		}
		return u.arco.NewSelect().
			Properties(structToMap(SelectProperties{
				Bordered:            true,
				UnmountOnExit:       true,
				Options:             options,
				Placeholder:         o.Default,
				Size:                "default",
				AutoAlignPopupWidth: true,
				Position:            "bottom",
				MountToBody:         true,
			}))
	}

	// TODO(xinxi.guo): implement validation
	return u.arco.NewInput().Properties(structToMap(InputProperties{
		Placeholder: o.Default,
		Size:        "default",
	}))
}

func (u UI) headerElements() []sunmao.BaseComponentBuilder {
	title := u.arco.NewText().Content(u.cli.Name)

	help := u.arco.NewButton().
		Properties(structToMap(ButtonProperties{
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
		Properties(structToMap(ButtonProperties{
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
