package ui

import (
	"CLI2UI/pkg/config"
	"fmt"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

func (u UI) buildPage() {
	f := u.cli.Form()

	cs := []sunmao.BaseComponentBuilder{
		u.layout(),
		u.helpModal(),
	}

	for _, c := range cs {
		u.arco.Component(c)
	}

	u.registerEvents(&f)
}

type UpdateSubcommandParams[T int | string] struct {
	Path            Path
	SubcommandIndex T
	Tabs            []TabProperties
}

func (u UI) registerEvents(f *config.Form) {
	u.r.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		p := toStruct[UpdateSubcommandParams[int]](m.Params)
		form := p.Path.traverseForm(f)
		form.Choice = p.Tabs[p.SubcommandIndex].Title
		clearForm(form)
		return nil
	})
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
			},
		})
}

func (u UI) commandStack(p Path, c config.Command) *sunmao.StackComponentBuilder {
	cs := []sunmao.BaseComponentBuilder{u.optionsInputForm(p, c)}

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

	for _, c := range c.Subcommands {
		tabs = append(tabs, TabProperties{
			Title:         c.Name,
			DestroyOnHide: true,
		})
	}

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

	// TODO(xinxi.guo): disable when len(options) == 0
	cb := u.arco.NewCheckbox().
		Id(p.optionsCheckboxId()).
		Properties(structToMap(CheckboxProperties{
			Options:              os,
			DefaultCheckedValues: required,
			Direction:            "horizontal",
		})).
		Style("content", `
		display: flex;
		align-items: flex-end;
		`)

	// TODO(xinxi.guo): notify the user when there is no option available
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
		inputs = append(inputs, u.optionInput(p, c.Name, f))
		os = append(os, CheckboxOptionProperties{
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
		os = append(os, CheckboxOptionProperties{
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

func (u UI) optionInput(p Path, cmd string, o config.FlagOrArg) sunmao.BaseComponentBuilder {
	return u.arco.NewFormControl().
		Id(p.optionValueInputId(o.Name)).
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
			ID:   p.optionValuesFormId(),
			Slot: "content",
		}, fmt.Sprintf("{{ %s.checkedValues.some(o => o === \"%s\") }}", p.optionsCheckboxId(), o.Name))
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
