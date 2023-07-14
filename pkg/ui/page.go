package ui

import (
	"CLI2UI/pkg/config"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

func (u UI) buildPage() {
	cs := []sunmao.BaseComponentBuilder{
		u.layout(),
	}

	for _, c := range cs {
		u.arco.Component(c)
	}
}

func (u UI) layout() sunmao.BaseComponentBuilder {
	return u.arco.NewStack().
		Properties(structToMap(StackProperties{
			Direction: "horizontal",
		})).
		Style("content", `
		user-select: none;
		height: 100vh;
		width: 100vw;
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				u.sidebar(),
				u.mainContent(),
			},
		})
}

func (u UI) mainContent() sunmao.BaseComponentBuilder {
	return u.arco.NewStack().
		Style("content", `
		display: grid;
		flex: 1;
		`).Children(map[string][]sunmao.BaseComponentBuilder{
		"content": {},
	})
}

func (u UI) sidebar() sunmao.BaseComponentBuilder {
	title := u.arco.NewText().
		Content(u.cli.Name).
		Style("content", `
		font-size: 1.25rem;
		font-weight: bold;
		margin: 0.125rem 0;
		min-width: 192px;
		`)

	s := u.arco.NewStack().
		Properties(structToMap(StackProperties{
			Direction: "vertical",
			Spacing:   6,
		})).
		Style("content", `
		padding: 0.75rem;
		color: #fff;
		overflow: auto;
		background-color: rgba(11, 21, 48, 0.9); 
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				title,
				u.commandMenu(),
			},
		})

	return s
}

func (u UI) commandMenu() sunmao.BaseComponentBuilder {
	return u.arco.NewTree().
		Properties(structToMap(
			TreeProperties{
				Data: u.menuItems(),
				Size: "large",
			},
		)).
		Style("content", `
		.arco-tree-node {
			color: #fff;
		}
		`)
}

func (u UI) menuItems() []TreeNodeProperties {
	return menuItems(u.cli.Command, []TreeNodeProperties{})
}

func menuItems(c config.Command, i []TreeNodeProperties) []TreeNodeProperties {
	for _, sc := range c.Subcommands {
		tnp := TreeNodeProperties{
			Title:      sc.DisplayName(),
			Key:        sc.Name,
			Children:   menuItems(sc, []TreeNodeProperties{}),
			Selectable: true,
		}
		i = append(i, tnp)
	}

	return i
}
