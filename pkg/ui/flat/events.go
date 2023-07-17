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

		return pathState.SetState(p.Path, &connId)
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

	u.Runtime.Handle("UpdateCheckedOptions", func(m *runtime.Message, connId int) error {
		s := ui.GetOrCreateSession(*u.FormTemplate, connId)
		p := ui.ToStruct[ui.UpdateCheckedOptionsParams[[]string]](m.Params)
		f := p.Path.TraverseForm(s.Form)
		ui.UpdateCheckedOptions(&f.Flags, p.CheckedValues)
		ui.UpdateCheckedOptions(&f.Args, p.CheckedValues)
		return nil
	})

	u.Runtime.Handle("EstablishedConnection", func(m *runtime.Message, connId int) error {
		ui.GetOrCreateSession(*u.FormTemplate, connId)
		return nil
	})
}
