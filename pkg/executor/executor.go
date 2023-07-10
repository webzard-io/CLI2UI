package executor

import (
	"github.com/go-cmd/cmd"
	"github.com/google/shlex"
)

type Executor struct {
	State   *ExecuteState
	stateCh chan *ExecuteState
	stopCh  chan struct{}
}

type ExecuteState struct {
	Error     error  `json:"error"`
	Done      bool   `json:"done"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	IsRunning bool   `json:"isRunning"`
}

func NewExecutor(stateCh chan *ExecuteState, stopCh chan struct{}) Executor {
	state := &ExecuteState{}

	return Executor{
		State:   state,
		stateCh: stateCh,
		stopCh:  stopCh,
	}
}

func (e *Executor) Run(script string) error {
	e.resetState()

	frags, err := shlex.Split(script)
	if err != nil {
		return err
	}

	c := cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, frags[0], frags[1:]...)
	statusCh := c.Start()

	go func() {
		for !e.State.Done && (c.Stdout != nil || c.Stderr != nil) {
			e.State.IsRunning = true
			select {
			case line, open := <-c.Stdout:
				if !open {
					c.Stdout = nil
					continue
				}
				e.State.Stdout = e.State.Stdout + line + "\r\n"
				e.stateCh <- e.State
			case line, open := <-c.Stderr:
				if !open {
					c.Stderr = nil
					continue
				}
				e.State.Stderr = e.State.Stderr + line + "\r\n"
				e.stateCh <- e.State
			}
		}
	}()

	go func() {
		<-e.stopCh
		c.Stop()
	}()

	go func() {
		finalStatus := <-statusCh
		e.State.Done = true
		e.State.Error = finalStatus.Error
		e.State.IsRunning = false
		e.stateCh <- e.State
	}()

	return nil
}

func (e *Executor) resetState() {
	e.State = &ExecuteState{
		Error:     nil,
		Done:      false,
		Stderr:    "",
		Stdout:    "",
		IsRunning: false,
	}
	e.stateCh <- e.State
}
