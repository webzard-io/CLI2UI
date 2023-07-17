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

func (u UI) registerEvents() {
	execState := u.Runtime.NewServerState("exec", nil)
	u.Arco.Component(execState.AsComponent())

	dryRunState := u.Runtime.NewServerState("dryRun", "")
	u.Arco.Component(dryRunState.AsComponent())

	formatState := u.Runtime.NewServerState("format", "")
	u.Arco.Component(formatState.AsComponent())

	u.Runtime.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.FormTemplate, connId)

		p := ui.ToStruct[UpdateSubcommandParams[int]](m.Params)
		form := p.Path.TraverseForm(s.Form)
		form.Choice = p.Values[p.SubcommandIndex]
		form.Subcommands[form.Choice].Clear()
		return nil
	})

	u.Runtime.Handle("UpdateOptionValue", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.FormTemplate, connId)

		p := ui.ToStruct[ui.UpdateOptionValueParams](m.Params)
		form := p.Path.TraverseForm(s.Form)

		_, ok := form.Args[p.OptionName]
		if ok {
			form.Args[p.OptionName].Value = p.Value
		} else {
			form.Flags[p.OptionName].Value = p.Value
		}

		return nil
	})

	u.Runtime.Handle("Heartbeat", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.FormTemplate, connId)
		s.HeartbeatCh <- struct{}{}
		return nil
	})

	u.Runtime.Handle("Run", func(m *runtime.Message, connId int) error {
		sess := ui.GetOrCreateSession(*u.FormTemplate, connId)

		script, f := u.CLI.Script(*sess.Form)
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
				err := u.Runtime.Ping(&connId, "Ping")

				// this fails when a WebSocket connection drops **loudly**
				if err != nil {
					sess.Exec.StopCh <- struct{}{}
					return
				}

				select {
				case <-sess.HeartbeatCh:
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

	u.Runtime.Handle("Stop", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.FormTemplate, connId)
		s.Exec.StopCh <- struct{}{}
		return nil
	})

	u.Runtime.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		ui.GetOrCreateSession(*u.FormTemplate, connId)
		return nil
	})

	u.Runtime.Handle("DryRun", func(m *runtime.Message, connId int) error {
		sess := ui.GetOrCreateSession(*u.FormTemplate, connId)
		s, _ := u.CLI.Script(*sess.Form)
		return dryRunState.SetState(s, &connId)
	})

	u.Runtime.Handle("UpdateCheckedOptions", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.FormTemplate, connId)
		p := ui.ToStruct[UpdateCheckedOptionsParams[[]string]](m.Params)
		f := p.Path.TraverseForm(s.Form)
		ui.UpdateCheckedOptions(&f.Flags, p.CheckedValues)
		ui.UpdateCheckedOptions(&f.Args, p.CheckedValues)
		return nil
	})
}
