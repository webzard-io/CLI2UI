package naive

import (
	"CLI2UI/pkg/ui"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
)

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

	formatState := u.r.NewServerState("format", "")
	u.arco.Component(formatState.AsComponent())

	u.r.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.fTpl, connId)

		p := ui.ToStruct[UpdateSubcommandParams[int]](m.Params)
		form := p.Path.TraverseForm(s.Form)
		form.Choice = p.Values[p.SubcommandIndex]
		form.Subcommands[form.Choice].Clear()
		return nil
	})

	u.r.Handle("UpdateOptionValue", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.fTpl, connId)

		p := ui.ToStruct[UpdateOptionValueParams](m.Params)
		form := p.Path.TraverseForm(s.Form)

		_, ok := form.Args[p.OptionName]
		if ok {
			form.Args[p.OptionName].Value = p.Value
		} else {
			form.Flags[p.OptionName].Value = p.Value
		}

		return nil
	})

	u.r.Handle("Heartbeat", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.fTpl, connId)
		s.HeatbeatCh <- struct{}{}
		return nil
	})

	u.r.Handle("Run", func(m *runtime.Message, connId int) error {
		sess := ui.GetOrCreateSession(*u.fTpl, connId)

		script, f := u.cli.Script(*sess.Form)
		formatState.SetState(f, &connId)

		finishedCh, err := sess.Exec.Run(script)
		if err != nil {
			log.Error(err)
			return err
		}

		go func() {
			for {
				select {
				case <-sess.Exec.StateCh:
					err := execState.SetState(sess.Exec.State, &connId)
					if err != nil {
						log.Error(err)
					}
				case <-finishedCh:
					err := execState.SetState(sess.Exec.State, &connId)
					if err != nil {
						log.Error(err)
					}
					return
				}
			}
		}()

		// force an update to the state so the run button is disabled
		sess.Exec.StateCh <- struct{}{}

		go func() {
			for sess.Exec.State.IsRunning {
				// TODO(xinxi.guo): this can be extended to send more useful messages
				err := u.r.Ping(&connId, "Ping")

				// this fails when a WebSocket connection drops **loudly**
				if err != nil {
					sess.Exec.StopCh <- struct{}{}
					return
				}

				select {
				case <-sess.HeatbeatCh:
				// this fails when a WebSocket connection drops **silently**
				case <-time.After(5 * time.Second):
					sess.Exec.StopCh <- struct{}{}
					return
				}

				time.Sleep(5 * time.Second)
			}
		}()

		return nil
	})

	u.r.Handle("Stop", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.fTpl, connId)
		s.Exec.StopCh <- struct{}{}
		return nil
	})

	u.r.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		ui.GetOrCreateSession(*u.fTpl, connId)
		return nil
	})

	u.r.Handle("DryRun", func(m *runtime.Message, connId int) error {
		sess := ui.GetOrCreateSession(*u.fTpl, connId)
		s, _ := u.cli.Script(*sess.Form)
		return dryRunState.SetState(s, &connId)
	})

	u.r.Handle("UpdateCheckedOptions", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.fTpl, connId)
		p := ui.ToStruct[UpdateCheckedOptionsParams[[]string]](m.Params)
		f := p.Path.TraverseForm(s.Form)
		updateCheckedOptions(&f.Flags, p.CheckedValues)
		updateCheckedOptions(&f.Args, p.CheckedValues)
		return nil
	})
}
