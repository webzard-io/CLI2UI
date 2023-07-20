package ui

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/executor"
)

var sessions = map[int]*session{}

type session struct {
	Form        *config.Form
	Exec        *executor.Executor
	HeartbeatCh chan struct{}
	cliIndex    int
}

func GetOrCreateSession(cliIndex int, templates []*config.Form, connId int) *session {
	s, ok := sessions[connId]
	if !ok {
		f := templates[cliIndex].Clone()
		hbCh := make(chan struct{})
		exec := executor.NewExecutor()

		s = &session{
			Form:        f,
			Exec:        &exec,
			HeartbeatCh: hbCh,
			cliIndex:    cliIndex,
		}
		sessions[connId] = s
	}

	if cliIndex != s.cliIndex {
		s.Form = templates[cliIndex].Clone()
		s.cliIndex = cliIndex
	}

	return s
}
