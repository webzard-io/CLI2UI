package ui

import (
	"CLI2UI/pkg/config"

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
			"header":  u.headerElements(),
			"content": {u.flagsAndArgs(u.cli.Command)},
		})
}

func (u UI) flagsAndArgs(c config.Command) sunmao.BaseComponentBuilder {
	options, requiredOptions := parseOptions(c)

	cb := u.c2u.NewCheckboxMenu().
		Properties(structToMap(CheckboxMenuProperties{
			Value:   requiredOptions,
			Text:    "Options",
			Options: options,
		}))

	s := u.arco.NewStack().
		Properties(structToMap(StackProperties{
			Direction: "vertical",
		})).
		Style("content", "margin-left: auto;").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {cb},
		})

	return s
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
		Style("content", "width: 80vw").
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {help},
			"footer":  {close},
		})

	return modal
}
