package ui

import "github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"

func (u UI) buildPage() {
	l := u.layout()
	u.arco.Component(l)
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
		Children(map[string][]sunmao.BaseComponentBuilder{
			"header":  u.headerElements(),
			"content": {},
		})
}

func (u UI) headerElements() []sunmao.BaseComponentBuilder {
	title := u.arco.NewText().Content(u.cli.Name)
	help := u.arco.NewButton().
		Properties(structToMap(ButtonProperties{
			Shape: "square",
			Text:  "Help",
			Type:  "default",
		}))

	return []sunmao.BaseComponentBuilder{title, help}
}
