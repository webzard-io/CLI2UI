package ui

import (
	"CLI2UI/pkg/executor"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
)

type UpdateSubcommandParams[T int | string] struct {
	Path            Path
	SubcommandIndex T
	Tabs            []TabProperties
}

type UpdateOptionValueParams struct {
	Path       Path
	OptionName string
	Value      any
}

// TODO(xinxi.guo): different connId/user will have different form/exec/state to use
func (u UI) registerEvents() {
	f := *u.fTpl

	u.r.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		p := toStruct[UpdateSubcommandParams[int]](m.Params)
		form := p.Path.traverseForm(&f)
		form.Choice = p.Tabs[p.SubcommandIndex].Title
		clearForm(form)
		return nil
	})

	u.r.Handle("UpdateOptionValue", func(m *runtime.Message, connId int) error {
		p := toStruct[UpdateOptionValueParams](m.Params)
		form := p.Path.traverseForm(&f)

		_, ok := form.Args[p.OptionName]
		if ok {
			form.Args[p.OptionName].Value = p.Value
		} else {
			form.Flags[p.OptionName].Value = p.Value
		}

		return nil
	})

	hbCh := make(chan struct{})
	u.r.Handle("Heartbeat", func(m *runtime.Message, connId int) error {
		hbCh <- struct{}{}
		return nil
	})

	stateCh := make(chan *executor.ExecuteState)
	stopCh := make(chan struct{})

	exec := executor.NewExecutor(stateCh, stopCh)
	execState := u.r.NewServerState("execState", exec.State)
	u.arco.Component(execState.AsComponent())

	u.r.Handle("Run", func(m *runtime.Message, connId int) error {
		s := u.cli.Script(f)

		go func() {
			for !exec.State.Done {
				s := <-stateCh
				err := execState.SetState(s, &connId)
				if err != nil {
					log.Error(err)
				}
			}
		}()

		go func() {
			for !exec.State.Done {
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

		err := exec.Run(s)
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	u.r.Handle("Stop", func(m *runtime.Message, connId int) error {
		stopCh <- struct{}{}
		return nil
	})
}
