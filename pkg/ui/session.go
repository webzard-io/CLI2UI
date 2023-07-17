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
}

func GetOrCreateSession(template config.Form, connId int) *session {
	s, ok := sessions[connId]
	if !ok {
		f := template.Clone()
		hbCh := make(chan struct{})
		exec := executor.NewExecutor()

		s = &session{
			Form:        f,
			Exec:        &exec,
			HeartbeatCh: hbCh,
		}
		sessions[connId] = s
	}
	return s
}
