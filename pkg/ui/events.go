package ui

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/executor"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
)

var sessions = map[int]*session{}

type session struct {
	f    *config.Form
	exec *executor.Executor
	hbCh chan struct{}
}

func (u UI) GetOrCreateSession(connId int) *session {
	s, ok := sessions[connId]
	if !ok {
		f := u.fTpl.Clone()
		hbCh := make(chan struct{})
		exec := executor.NewExecutor()

		s = &session{
			f:    f,
			exec: &exec,
			hbCh: hbCh,
		}
		sessions[connId] = s
	}
	return s
}

type UpdateSubcommandParams[T int | string] struct {
	Path            Path
	SubcommandIndex T
	Values          []string
}

type UpdateCheckedOptionsParams[T []string | string] struct {
	Path          Path
	CheckedValues T
}

type UpdateOptionValueParams struct {
	Path       Path
	OptionName string
	Value      any
}

func (u UI) registerEvents() {
	execState := u.r.NewServerState("exec", nil)
	u.arco.Component(execState.AsComponent())

	dryRunState := u.r.NewServerState("dryRun", "")
	u.arco.Component(dryRunState.AsComponent())

	u.r.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		s := u.GetOrCreateSession(connId)

		p := toStruct[UpdateSubcommandParams[int]](m.Params)
		form := p.Path.traverseForm(s.f)
		form.Choice = p.Values[p.SubcommandIndex]
		form.Subcommands[form.Choice].Clear()
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

		finishedCh, err := sess.exec.Run(u.cli.Script(*sess.f))
		if err != nil {
			log.Error(err)
			return err
		}

		go func() {
			for {
				select {
				case <-sess.exec.StateCh:
					err := execState.SetState(sess.exec.State, &connId)
					if err != nil {
						log.Error(err)
					}
				case <-finishedCh:
					err := execState.SetState(sess.exec.State, &connId)
					if err != nil {
						log.Error(err)
					}
					return
				}
			}
		}()

		go func() {
			for sess.exec.State.IsRunning {
				// TODO(xinxi.guo): this can be extended to send more useful messages
				err := u.r.Ping(&connId, "Ping")

				// this fails when a WebSocket connection drops **loudly**
				if err != nil {
					sess.exec.StopCh <- struct{}{}
					return
				}

				select {
				case <-sess.hbCh:
				// this fails when a WebSocket connection drops **silently**
				case <-time.After(5 * time.Second):
					sess.exec.StopCh <- struct{}{}
					return
				}

				time.Sleep(5 * time.Second)
			}
		}()

		return nil
	})

	u.r.Handle("Stop", func(m *runtime.Message, connId int) error {
		s := u.GetOrCreateSession(connId)
		s.exec.StopCh <- struct{}{}
		return nil
	})

	u.r.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		u.GetOrCreateSession(connId)
		return nil
	})

	u.r.Handle("DryRun", func(m *runtime.Message, connId int) error {
		sess := u.GetOrCreateSession(connId)
		s := u.cli.Script(*sess.f)
		return dryRunState.SetState(s, &connId)
	})

	u.r.Handle("UpdateCheckedOptions", func(m *runtime.Message, connId int) error {
		s := u.GetOrCreateSession(connId)
		p := toStruct[UpdateCheckedOptionsParams[[]string]](m.Params)
		f := p.Path.traverseForm(s.f)
		updateCheckedOptions(&f.Flags, p.CheckedValues)
		updateCheckedOptions(&f.Args, p.CheckedValues)
		return nil
	})
}
