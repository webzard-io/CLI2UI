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
	p := prefix("", u.cli.Command.Name)

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
			},
		})
}

func (u UI) commandStack(p Prefix, c config.Command) *sunmao.StackComponentBuilder {
	cs := []sunmao.BaseComponentBuilder{u.optionsInputForm(p, c)}

	if len(c.Subcommands) > 0 {
		cs = append(cs, u.subcommandsTab(p, c))
		for i, c := range c.Subcommands {
			pre := prefix(p, c.Name)
			s := u.commandStack(pre, c).Slot(sunmao.Container{
				ID:   commandStackId(p),
				Slot: "content",
			}, fmt.Sprintf("{{ %s.activeTab === %d }}", subcommandTabsId(p), i))
			cs = append(cs, s)
		}
	}

	return u.arco.NewStack().
		Id(commandStackId(p)).
		Properties(structToMap(StackProperties{
			Direction: "vertical",
			Spacing:   6,
		})).
		Style("content", "width: 100%;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": cs,
		})
}

func (u UI) subcommandsTab(p Prefix, c config.Command) sunmao.BaseComponentBuilder {
	tabs := []TabProperties{}

	for _, c := range c.Subcommands {
		tabs = append(tabs, TabProperties{
			Title:         c.Name,
			DestroyOnHide: true,
		})
	}

	return u.arco.NewTabs().
		Id(subcommandTabsId(p)).
		Properties(structToMap(TabsProperties{
			Type:        "line",
			TabPosition: "top",
			Size:        "default",
			Tabs:        tabs,
		})).
		Style("content", "width: 100%;")
}

func (u UI) optionsInputForm(p Prefix, c config.Command) sunmao.BaseComponentBuilder {
	os, required, inputs := u.parseOptions(p, c)

	// TODO(xinxi.guo): disable when len(options) == 0
	cb := u.c2u.NewCheckboxMenu().
		Id(optionsCheckboxId(p)).
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
		Id(optionValuesFormId(p)).
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

func (u UI) parseOptions(p Prefix, c config.Command) ([]CheckboxOption, []string, []sunmao.BaseComponentBuilder) {
	os := []CheckboxOption{}
	required := []string{}
	inputs := []sunmao.BaseComponentBuilder{}

	for _, f := range c.Flags {
		inputs = append(inputs, u.optionInput(p, c.Name, f))
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
		inputs = append(inputs, u.optionInput(p, c.Name, a))
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

func (u UI) optionInput(p Prefix, cmd string, o config.FlagOrArg) sunmao.BaseComponentBuilder {
	return u.arco.NewFormControl().
		Id(optionValueInputId(p, o.Name)).
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
			ID:   optionValuesFormId(p),
			Slot: "content",
		}, fmt.Sprintf("{{ %s.value.some(o => o === \"%s\") }}", optionsCheckboxId(p), o.Name))
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
