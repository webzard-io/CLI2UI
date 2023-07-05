package ui

import (
	"CLI2UI/pkg/config"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/arco"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

func NewUI(c config.CLI) *UI {
	r := runtime.New("ui", "patch")
	app := sunmao.NewApp()
	arco := arco.NewArcoApp(app)
	c2u := NewCLI2UIApp(app)
	fTpl := c.Form()

	return &UI{
		r:    r,
		arco: arco,
		c2u:  c2u,
		cli:  &c,
		fTpl: &fTpl,
	}
}

func (u UI) Run() error {
	u.buildPage()
	u.registerEvents()

	err := u.r.LoadApp(u.arco.AppBuilder)
	if err != nil {
		return err
	}

	u.r.Run()
	return nil
}

type UI struct {
	r    *runtime.Runtime
	arco *arco.ArcoAppBuilder
	c2u  *CLI2UIAppBuilder
	cli  *config.CLI
	fTpl *config.Form
}
