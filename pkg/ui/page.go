package ui

import (
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
	return u.arco.NewLayout().
		Properties(structToMap(LayoutProperties{
			ShowSideBar: true,
		})).
		Style("layout", `
		border: 0;
		`).
		Style("sidebar", `
		height: 100vh;
		overflow: auto;
		background-color: rgba(11, 21, 48, 0.9); 
		`).
		Style("content", `
		height: 100vh;
		overflow: auto;
		`)
}
