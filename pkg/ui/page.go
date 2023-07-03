package ui

import "github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"

func (u UI) buildPage() {

}

func (u UI) layout() sunmao.BaseComponentBuilder {
	return u.arco.NewLayout().Children(map[string][]sunmao.BaseComponentBuilder{})
}
