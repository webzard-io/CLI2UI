package ui

import (
	"CLI2UI/pkg/config"
	"CLI2UI/pkg/executor"

	"github.com/google/uuid"
)

var ServerSignature = uuid.NewString()
var sessions = map[string]*session{}
var connIdToClientId = map[int]string{}

type session struct {
	Form          *config.Form
	Exec          *executor.Executor
	HeartbeatCh   chan struct{}
	cliIndex      int
	CurrentConnId *int
}

func UpdateConnIdToClientId(connId int, clientId string) {
	connIdToClientId[connId] = clientId
}

func GetOrCreateSession(cliIndex int, templates []*config.Form, connId int) *session {
	sId := connIdToClientId[connId]
	s, ok := sessions[sId]
	if !ok {
		f := templates[cliIndex].Clone()
		hbCh := make(chan struct{})
		exec := executor.NewExecutor()

		s = &session{
			Form:          f,
			Exec:          &exec,
			HeartbeatCh:   hbCh,
			cliIndex:      cliIndex,
			CurrentConnId: &connId,
		}
		sessions[sId] = s
	}

	if cliIndex != s.cliIndex {
		s.Form = templates[cliIndex].Clone()
		s.cliIndex = cliIndex
	}

	return s
}
