package ui

import (
	"CLI2UI/pkg/executor"
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/arco"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type UI struct {
	r *runtime.Runtime
	b *arco.ArcoAppBuilder
}

func NewUI(stopCh chan struct{}) *UI {
	r := runtime.New("ui", "patch")
	app := sunmao.NewApp()
	b := arco.NewArcoApp(app)

	return &UI{
		r: r,
		b: b,
	}
}

func (u UI) Run() error {
	u.buildUI()

	err := u.r.LoadApp(u.b.AppBuilder)
	if err != nil {
		return err
	}

	u.r.Run()
	return nil
}

func (u UI) buildUI() {
	b := u.b

	stopCh := make(chan struct{})

	stateCh := make(chan *executor.ExecuteState)

	e := executor.NewExecutor(stateCh, stopCh)
	eState := u.r.NewServerState("exec", e.State)

	b.Component(eState.AsComponent())

	u.r.Handle("run", func(m *runtime.Message, connId int) error {
		command := ""
		b, _ := json.Marshal(m.Params)
		_ = json.Unmarshal(b, &command)

		go func() {
			for {
				select {
				case s := <-stateCh:
					err := eState.SetState(s, &connId)
					if err != nil {
						log.Error(err)
					}
					break
				}
			}
		}()
		err := e.Run(command)
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	u.r.Handle("stop", func(m *runtime.Message, connId int) error {
		stopCh <- struct{}{}
		return nil
	})

	pingRaw := CLIJson{
		Name: "ping",
		Help: `usage: ping [-AaDdfnoQqRrv] [-c count] [-G sweepmaxsize]
				[-g sweepminsize] [-h sweepincrsize] [-i wait]
				[-l preload] [-M mask | time] [-m ttl] [-p pattern]
				[-S src_addr] [-s packetsize] [-t timeout][-W waittime]
				[-z tos] host`,
		Commands: []Command{
			{
				Name: "ping",
				Flags: []FlagOrArg{
					{
						Name:     "count",
						Required: false,
						Type:     FlagArgTypeNumber,
						Short:    "c",
					},
				},
				Args: []FlagOrArg{
					{
						Name:     "host",
						Required: true,
						Type:     FlagArgTypeString,
					},
				},
			},
		},
	}

	componentSchemas := genSchemaComponents(pingRaw)

	for _, componentSchema := range componentSchemas {
		b.RawComponent(componentSchema)
	}
}
