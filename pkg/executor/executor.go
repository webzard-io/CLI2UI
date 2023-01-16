package executor

import (
	"github.com/go-cmd/cmd"
)

type Executor struct {
	State   *ExecuteState
	stateCh chan *ExecuteState
	stopCh  chan struct{}
}

type ExecuteState struct {
	Error  error  `json:"error"`
	Done   bool   `json:"done"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

func NewExecutor(stateCh chan *ExecuteState, stopCh chan struct{}) Executor {
	state := &ExecuteState{
		Error:  nil,
		Done:   false,
		Stderr: "",
		Stdout: "",
	}

	return Executor{
		State:   state,
		stateCh: stateCh,
		stopCh:  stopCh,
	}
}

func (e *Executor) Run(name string, args ...string) error {
	c := cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, name, args...)
	statusCh := c.Start()

	go func() {
		for c.Stdout != nil || c.Stderr != nil {
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
		for {
			select {
			case <-e.stopCh:
				c.Stop()
				break
			}
		}
	}()

	go func() {
		select {
		case finalStatus := <-statusCh:
			e.State.Done = true
			e.State.Error = finalStatus.Error
			e.stateCh <- e.State
			break
		}
	}()

	return nil
}
