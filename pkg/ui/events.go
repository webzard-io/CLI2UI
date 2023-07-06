package ui

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/executor"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
)

var sessions = map[int]*session{}

type session struct {
	f         *config.Form
	exec      *executor.Executor
	execState *runtime.ServerState
	stateCh   chan *executor.ExecuteState
	stopCh    chan struct{}
	hbCh      chan struct{}
	connId    int
}

func (u UI) GetOrCreateSession(connId int) *session {
	s, ok := sessions[connId]
	if !ok {
		f := new(config.Form)
		*f = *u.fTpl
		fmt.Println(*f)
		stopCh := make(chan struct{})
		stateCh := make(chan *executor.ExecuteState)
		hbCh := make(chan struct{})
		exec := executor.NewExecutor(stateCh, stopCh, connId)
		execState := u.r.NewServerState(executorId(connId), exec.State)
		u.arco.Component(execState.AsComponent())

		s := &session{
			f:         f,
			exec:      &exec,
			execState: execState,
			stateCh:   stateCh,
			stopCh:    stopCh,
			hbCh:      hbCh,
			connId:    connId,
		}
		sessions[connId] = s
	}
	return s
}

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
	u.r.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		s := u.GetOrCreateSession(connId)

		p := toStruct[UpdateSubcommandParams[int]](m.Params)
		form := p.Path.traverseForm(s.f)
		form.Choice = p.Tabs[p.SubcommandIndex].Title
		clearForm(form)
		return nil
	})

	u.r.Handle("UpdateOptionValue", func(m *runtime.Message, connId int) error {
		s := u.GetOrCreateSession(connId)

		p := toStruct[UpdateOptionValueParams](m.Params)
		form := p.Path.traverseForm(s.f)

		_, ok := form.Args[p.OptionName]
		if ok {
			form.Args[p.OptionName].Value = p.Value
		} else {
			form.Flags[p.OptionName].Value = p.Value
		}

		return nil
	})

	u.r.Handle("Heartbeat", func(m *runtime.Message, connId int) error {
		s := u.GetOrCreateSession(connId)
		s.hbCh <- struct{}{}
		return nil
	})

	u.r.Handle("Run", func(m *runtime.Message, connId int) error {
		sess := u.GetOrCreateSession(connId)

		go func() {
			for !sess.exec.State.Done {
				s := <-sess.stateCh
				err := sess.execState.SetState(s, &connId)
				if err != nil {
					log.Error(err)
				}
			}
		}()

		go func() {
			for !sess.exec.State.Done {
				// TODO(xinxi.guo): this can be extended to send more useful messages
				err := u.r.Ping(&connId, "Ping")

				// this fails when a WebSocket connection drops **loudly**
				if err != nil {
					sess.stopCh <- struct{}{}
					return
				}

				select {
				case <-sess.hbCh:
				// this fails when a WebSocket connection drops **silently**
				case <-time.After(5 * time.Second):
					sess.stopCh <- struct{}{}
					return
				}

				time.Sleep(5 * time.Second)
			}
		}()

		err := sess.exec.Run(u.cli.Script(*sess.f))
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})

	u.r.Handle("Stop", func(m *runtime.Message, connId int) error {
		s := u.GetOrCreateSession(connId)
		s.stopCh <- struct{}{}
		return nil
	})

	u.r.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		fmt.Println(connId)
		u.GetOrCreateSession(connId)
		fmt.Printf("Connection ID: %d\n", connId)
		return nil
	})
}
