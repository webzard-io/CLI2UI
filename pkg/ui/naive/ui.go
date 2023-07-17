package naive

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
)

func NewUI(c config.CLI) (*UI, error) {
	base, err := ui.NewUI(c)
	return &UI{base}, err
}

func (u UI) Run() error {
	u.buildPage()
	u.registerEvents()

	err := u.Runtime.LoadApp(u.Arco.AppBuilder)
	if err != nil {
		return err
	}

	u.Runtime.Run()
	return nil
}

type UI struct {
	*ui.BaseUI
}
