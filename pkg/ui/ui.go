package ui

import (
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type UI struct {
	r    *runtime.Runtime
	b    *sunmao.AppBuilder
	arco *sunmao.ArcoAppBuilder
}

func NewUI(stopCh chan struct{}) *UI {
	r := runtime.New("ui", "patch")
	b := sunmao.NewApp()
	arco := sunmao.NewArcoApp()

	return &UI{
		r:    r,
		b:    b,
		arco: arco,
	}
}

func (u UI) Run() error {
	err := u.r.LoadApp(u.b)
	if err != nil {
		return err
	}

	u.r.Run()
	return nil
}
