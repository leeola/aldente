package local

import "os/exec"

type Command struct {
	name        string
	machineName string
	C           *exec.Cmd
}

func (c Command) Name() string {
	return c.name
}

func (c Command) MachineName() string {
	return c.machineName
}

func (c Command) Start() error {
	return c.C.Start()
}

func (c Command) Wait() error {
	return c.C.Wait()
}
