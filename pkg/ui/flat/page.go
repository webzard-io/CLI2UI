package flat

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"

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
		Properties(ui.StructToMap(ui.StackProperties{
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
		Properties(ui.StructToMap(ui.StackProperties{
			Direction: "horizontal",
		})).
		Style("content", `
		flex: 1;
		background-color: rgb(241, 245, 249);
		gap: 0.5rem;
		padding: 0.5rem;
		`).Children(map[string][]sunmao.BaseComponentBuilder{
		"content": {
			u.optionSection(),
			u.outputSection(),
		},
	})
}

func (u UI) optionSection() sunmao.BaseComponentBuilder {
	return u.arco.NewStack().
		Style("content", `
		flex: 2;
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				u.arco.NewStack().
					Style("content", `
					flex: 1;
					background-color: white;
					border-radius: 0.5rem;
					padding: 0.5rem;
				`),
			},
		})
}

func (u UI) outputSection() sunmao.BaseComponentBuilder {
	stdoutCard := u.arco.NewStack().
		Properties(ui.StructToMap(ui.StackProperties{
			Direction: "vertical",
		})).
		Style("content", `
		flex: 1.5;
		background-color: white;
		border-radius: 0.5rem;
		padding: 0.5rem;
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				u.arco.NewText().
					Style("content", `
					font-size: 1.25rem;
					font-weight: bold;
					`).
					Content("Standard Output"),
			},
		})

	stderrCard := u.arco.NewStack().
		Properties(ui.StructToMap(ui.StackProperties{
			Direction: "vertical",
		})).
		Style("content", `
		flex: 1;
		background-color: white;
		border-radius: 0.5rem;
		padding: 0.5rem;
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				u.arco.NewText().
					Style("content", `
					font-size: 1.25rem;
					font-weight: bold;
					`).
					Content("Standard Error"),
			},
		})

	return u.arco.NewStack().
		Properties(ui.StructToMap(ui.StackProperties{
			Direction: "vertical",
		})).
		Style("content", `
		flex: 1;
		gap: 0.5rem;
		`).
		Children(map[string][]sunmao.BaseComponentBuilder{
			"content": {
				stdoutCard,
				stderrCard,
			},
		})
}

func (u UI) sidebar() sunmao.BaseComponentBuilder {
	title := u.arco.NewText().
		Content(u.cli.Name).
		Style("content", `
		font-size: 1.25rem;
		font-weight: bold;
		margin: 0.125rem 0;
		`)

	s := u.arco.NewStack().
		Properties(ui.StructToMap(ui.StackProperties{
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
		Id("SubcommandMenuTree").
		Properties(ui.StructToMap(
			ui.TreeProperties{
				Data: u.menuItems(),
				Size: "large",
			},
		)).
		Style("content", `
		.arco-tree-node {
			color: #fff;
		}
		`).
		Event([]sunmao.EventHandler{
			{
				Type:        "onSelect",
				ComponentId: "$utils",
				Method: sunmao.EventMethod{
					Name: "binding/v1/UpdateSubcommand",
					Parameters: UpdateSubcommandParams[string]{
						Path:       "{{ SubcommandMenuTree.selectedNodes[0].path }}",
						Subcommand: "{{ SubcommandMenuTree.selectedKeys[0] }}",
					},
				},
			},
		})
}

func (u UI) menuItems() []ui.TreeNodeProperties {
	return menuItems(u.cli.Command, []ui.TreeNodeProperties{})
}

func menuItems(c config.Command, i []ui.TreeNodeProperties) []ui.TreeNodeProperties {
	for _, sc := range c.Subcommands {
		tnp := ui.TreeNodeProperties{
			Title:      sc.DisplayName(),
			Key:        sc.Name,
			Children:   menuItems(sc, []ui.TreeNodeProperties{}),
			Selectable: true,
		}
		i = append(i, tnp)
	}

	return i
}
