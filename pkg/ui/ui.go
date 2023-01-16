package ui

import (
	"CLI2UI/pkg/executor"
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/sunmao"
)

type UI struct {
	r *runtime.Runtime
	b *sunmao.ArcoAppBuilder
}

func NewUI(stopCh chan struct{}) *UI {
	r := runtime.New("ui", "patch")
	b := sunmao.NewArcoApp()

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

	b.Component(b.NewStack().Children(map[string][]sunmao.BaseComponentBuilder{
		"content": {
			b.NewText().Content("CLI2UI"),
			b.NewComponent().Id("command_input").
				Type("arco/v1/input").
				Style("input", `
					width: 600px;
				`),
			b.NewComponent().
				Type("arco/v1/button").
				Properties(map[string]interface{}{
					"long": false,
					"type": "primary",
					"text": "run",
				}).
				Style("content", `width: 100px;`).
				Trait(b.NewTrait().Type("core/v1/event").Properties(map[string]interface{}{
					"handlers": []map[string]interface{}{
						{
							"type":        "onClick",
							"componentId": "$utils",
							"method": map[string]interface{}{
								"name":       fmt.Sprintf("binding/v1/%v", "run"),
								"parameters": "{{ command_input.value.split(' ') }}",
							},
						},
					},
				})),
			b.NewComponent().
				Type("arco/v1/button").
				Properties(map[string]interface{}{
					"long": false,
					"type": "default",
					"text": "stop",
				}).
				Style("content", `width: 100px;`).Trait(b.NewTrait().Type("core/v1/event").Properties(map[string]interface{}{
				"handlers": []map[string]interface{}{
					{
						"type":        "onClick",
						"componentId": "$utils",
						"method": map[string]interface{}{
							"name":       fmt.Sprintf("binding/v1/%v", "stop"),
							"parameters": "",
						},
					},
				},
			})),
			b.NewText().Content(`done: {{ exec.state.done }}`),
			b.NewText().Content(`error: {{ exec.state.error ? JSON.stringify(exec.state.error) : "-" }}`),
			b.NewText().Content(`stdout:`),
			b.NewText().Content(`{{ exec.state.stdout }}`).Style("content", "white-space:pre;"),
			b.NewText().Content(`stderr:`),
			b.NewText().Content(`{{ exec.state.stderr }}`).Style("content", "white-space:pre;"),
		},
	}).Properties(map[string]interface{}{
		"direction": "vertical",
	}).Style("content", `
		width: 100%;
		padding: 1em;
		margin-bottom: .5em;
	`))

	u.r.Handle("run", func(m *runtime.Message, connId int) error {
		command := make([]string, 0)
		b, _ := json.Marshal(m.Params)
		_ = json.Unmarshal(b, &command)

		go func() {
			for {
				select {
				case s := <-stateCh:
					eState.SetState(s, &connId)
					break
				}
			}
		}()

		err := e.Run(command[0], command[1:]...)
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
}
