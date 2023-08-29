package flat

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/ui"
	"fmt"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
)

type UpdateSubcommandParams[T string | ui.Path] struct {
	Path       T      `json:"myPath"`
	Subcommand string `json:"subcommand"`
}

func (u UI) registerEvents() {
	execState := u.Runtime.NewServerState("exec", nil)
	u.Arco.Component(execState.AsComponent())

	dryRunState := u.Runtime.NewServerState("dryRun", "")
	u.Arco.Component(dryRunState.AsComponent())

	formatState := u.Runtime.NewServerState("format", "")
	u.Arco.Component(formatState.AsComponent())

	pathState := u.Runtime.NewServerState("path", []string{})
	u.Arco.Component(pathState.AsComponent())

	u.Runtime.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(0, u.FormTemplates, connId)
		p := ui.ToStruct[UpdateSubcommandParams[ui.Path]](m.Params)

		form := p.Path.TraverseFormWithCallback(s.Form, func(s string, f *config.Form) {
			f.Choice = s
		})

		form.Choice = p.Subcommand
		form.Clear()

		return pathState.SetState(p.Path, &connId)
	})

	u.Runtime.Handle("UpdateOptionValue", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(0, u.FormTemplates, connId)

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
		s := ui.GetOrCreateSession(0, u.FormTemplates, connId)
		s.HeartbeatCh <- struct{}{}
		return nil
	})

	u.Runtime.Handle("Run", func(m *runtime.Message, connId int) error {
		sess := ui.GetOrCreateSession(0, u.FormTemplates, connId)

		script, f := u.CLIs[0].Script(*sess.Form)
		formatState.SetState(f, &connId)

		stdoutFile, err := os.OpenFile(fmt.Sprintf("stdout-%d.log", connId), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		stderrFile, err := os.OpenFile(fmt.Sprintf("stdout-%d.err", connId), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		finishedCh, err := sess.Exec.Run(nil, stdoutFile, stderrFile, script)
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
					stdoutFile.Close()
					stderrFile.Close()
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
		s := ui.GetOrCreateSession(0, u.FormTemplates, connId)
		s.Exec.StopCh <- struct{}{}
		return nil
	})

	u.Runtime.Handle("DryRun", func(m *runtime.Message, connId int) error {
		sess := ui.GetOrCreateSession(0, u.FormTemplates, connId)
		s, _ := u.CLIs[0].Script(*sess.Form)
		return dryRunState.SetState(s, &connId)
	})

	u.Runtime.Handle("UpdateCheckedOptions", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(0, u.FormTemplates, connId)
		p := ui.ToStruct[ui.UpdateCheckedOptionsParams[[]string]](m.Params)
		f := p.Path.TraverseForm(s.Form)
		ui.UpdateCheckedOptions(&f.Flags, p.CheckedValues)
		ui.UpdateCheckedOptions(&f.Args, p.CheckedValues)
		return nil
	})

	u.Runtime.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		ui.GetOrCreateSession(0, u.FormTemplates, connId)
		return nil
	})
}
