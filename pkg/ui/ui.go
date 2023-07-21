package ui

import (
	"CLI2UI/pkg/config"
	client "CLI2UI/ui"
	"errors"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/arco"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type UI interface {
	Run() error
}

func NewUI(c ...config.CLI) (*BaseUI, error) {
	if client.Error != nil {
		return nil, errors.New("failed to load prebuilt UI")
	}

	r := runtime.New(client.FS, "patch")
	app := sunmao.NewApp()
	arco := arco.NewArcoApp(app)
	c2u := NewCLI2UIApp(app)

	tpls := []*config.Form{}

	for _, cli := range c {
		f := cli.Form()
		tpls = append(tpls, &f)
	}

	return &BaseUI{
		Runtime:       r,
		Arco:          arco,
		C2U:           c2u,
		CLIs:          c,
		FormTemplates: tpls,
	}, nil
}

type BaseUI struct {
	Runtime       *runtime.Runtime
	Arco          *arco.ArcoAppBuilder
	C2U           *CLI2UIAppBuilder
	CLIs          []config.CLI
	FormTemplates []*config.Form
}

func (BaseUI) Run() error {
	return errors.New("cannot run BaseUI")
}
