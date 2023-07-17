package flat

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
	client "CLI2UI/ui"
	"errors"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/arco"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

func NewUI(c config.CLI) (*UI, error) {
	if client.Error != nil {
		return nil, errors.New("failed to load prebuilt UI")
	}

	r := runtime.New(client.FS, "patch")
	app := sunmao.NewApp()
	arco := arco.NewArcoApp(app)
	c2u := ui.NewCLI2UIApp(app)
	fTpl := c.Form()

	return &UI{
		r:    r,
		arco: arco,
		c2u:  c2u,
		cli:  &c,
		fTpl: &fTpl,
	}, nil
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
	c2u  *ui.CLI2UIAppBuilder
	cli  *config.CLI
	fTpl *config.Form
}
