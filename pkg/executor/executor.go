package executor

import (
	"github.com/go-cmd/cmd"
	"github.com/google/shlex"
)

type Executor struct {
	State   *state        `json:"state"`
	StateCh chan struct{} `json:"stateCh"`
	StopCh  chan struct{} `json:"stopCh"`
}

type state struct {
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	Error     error  `json:"error"`
	IsRunning bool   `json:"isRunning"`
}

func (e *Executor) Run(script string) (chan struct{}, error) {
	e.resetState()
	e.State.IsRunning = true

	frags, err := shlex.Split(script)
	if err != nil {
		return nil, err
	}

	c := cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, frags[0], frags[1:]...)

	finStatusCh := c.Start()

	stdioStatusCh := make(chan struct{})
	go func() {
		defer close(stdioStatusCh)
		for c.Stdout != nil || c.Stderr != nil {
			select {
			case line, open := <-c.Stdout:
				if !open {
					c.Stdout = nil
					continue
				}
				e.State.Stdout = e.State.Stdout + line + "\r\n"
			case line, open := <-c.Stderr:
				if !open {
					c.Stderr = nil
					continue
				}
				e.State.Stderr = e.State.Stderr + line + "\r\n"
			}
			e.StateCh <- struct{}{}
		}
	}()

	go func() {
		select {
		case <-e.StopCh:
			c.Stop()
		case <-stdioStatusCh:
		}
	}()

	finishedCh := make(chan struct{})

	go func() {
		finalStatus := <-finStatusCh
		e.State.Error = finalStatus.Error
		<-stdioStatusCh
		e.State.IsRunning = false
		e.StateCh <- struct{}{}
		close(finishedCh)
	}()

	return finishedCh, nil
}

func NewExecutor() Executor {
	return Executor{
		StateCh: make(chan struct{}),
		StopCh:  make(chan struct{}),
	}
}

func (e *Executor) resetState() {
	e.State = &state{
		Stdout:    "",
		Stderr:    "",
		Error:     nil,
		IsRunning: false,
	}
}
