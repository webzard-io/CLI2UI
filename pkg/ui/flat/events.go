package flat

import (
	"CLI2UI/pkg/ui"

	"github.com/yuyz0112/sunmao-ui-go-binding/pkg/runtime"
)

type UpdateSubcommandParams[T string | ui.Path] struct {
	Path       T
	Subcommand string
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
		s := ui.GetOrCreateSession(*u.FormTemplate, connId)
		p := ui.ToStruct[UpdateSubcommandParams[ui.Path]](m.Params)
		form := p.Path.TraverseForm(s.Form)
		form.Choice = p.Subcommand
		form.Clear()
		return nil
	})

	u.Runtime.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		ui.GetOrCreateSession(*u.FormTemplate, connId)
		return nil
	})
}
