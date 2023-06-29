package ui

import (
	"CLI2UI/pkg/executor"
	"encoding/json"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/arco"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type UI struct {
	r  *runtime.Runtime
	b  *arco.ArcoAppBuilder
	bb *CLI2UIAppBuilder
}

func NewUI(stopCh chan struct{}) *UI {
	r := runtime.New("ui", "patch")
	app := sunmao.NewApp()
	b := arco.NewArcoApp(app)
	bb := NewCLI2UIApp(app)

	return &UI{
		r:  r,
		b:  b,
		bb: bb,
	}
}

func (u *UI) Run() error {
	u.buildUI()

	err := u.r.LoadApp(u.b.AppBuilder)
	if err != nil {
		return err
	}

	u.r.Run()
	return nil
}

func (u *UI) buildUI() {
	stopCh := make(chan struct{})

	stateCh := make(chan *executor.ExecuteState)

	e := executor.NewExecutor(stateCh, stopCh)
	eState := u.r.NewServerState("exec", e.State)

	u.b.Component(eState.AsComponent())

	hbCh := make(chan struct{})
	u.r.Handle("Heartbeat", func(m *runtime.Message, connId int) error {
		hbCh <- struct{}{}
		return nil
	})

	u.r.Handle("run", func(m *runtime.Message, connId int) error {
		command := ""
		b, _ := json.Marshal(m.Params)
		_ = json.Unmarshal(b, &command)

		go func() {
			for !e.State.Done {
				s := <-stateCh
				err := eState.SetState(s, &connId)
				if err != nil {
					log.Error(err)
				}
			}
		}()

		go func() {
			for !e.State.Done {
				// TODO(xinxi.guo): this can be extended to send more useful messages
				err := u.r.Ping(&connId, "Ping")

				// this fails when a WebSocket connection drops **loudly**
				if err != nil {
					stopCh <- struct{}{}
					return
				}

				select {
				case <-hbCh:
				// this fails when a WebSocket connection drops **silently**
				case <-time.After(5 * time.Second):
					stopCh <- struct{}{}
					return
				}

				time.Sleep(5 * time.Second)
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

	u.genSchemaComponents(pingRaw)
}
