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
	execState := u.r.NewServerState("exec", nil)
	u.arco.Component(execState.AsComponent())

	dryRunState := u.r.NewServerState("dryRun", "")
	u.arco.Component(dryRunState.AsComponent())

	formatState := u.r.NewServerState("format", "")
	u.arco.Component(formatState.AsComponent())

	u.r.Handle("UpdateSubcommand", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.fTpl, connId)
		p := ui.ToStruct[UpdateSubcommandParams[ui.Path]](m.Params)
		form := p.Path.TraverseForm(s.Form)
		form.Choice = p.Subcommand
		form.Clear()
		return nil
	})

	u.r.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		ui.GetOrCreateSession(*u.fTpl, connId)
		return nil
	})
}
