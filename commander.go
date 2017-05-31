package aldente

import (
	"sync"

	"github.com/leeola/errors"
)

type Command interface {
	// CommandName returns the name the running command.
	CommandName() string

	// MachineName returns the machine name that is running the command.
	MachineName() string

	// Output returns a channel to monitor the progress of a command.
	//
	// If an error is encountered, it can be found Wait() return value.
	Output() <-chan CommandOutput

	// Wait for the entire Command process to be done.
	Wait() error
}

// CommandOutput contains a state and message sent during the provisioning.
type CommandOutput struct {
	// CommandName of the command that is running.
	CommandName string `json:"commandName"`

	// MachineName returns the machine name that is being provisioned.
	MachineName string `json:"machineName"`

	// Message is an output, likely a text line, from the given machine.
	Message string `json:"message"`
}

type Commands []Command

type CommandsError struct {
	CommandName string
	MachineName string
	Err         error
}

func (cmds Commands) Output() <-chan CommandOutput {
	c := make(chan CommandOutput, 10)
	w := &sync.WaitGroup{}
	w.Add(len(ps))

	for _, cmd := range cmds {
		go func(c chan CommandOutput, cmd Command, w *sync.WaitGroup) {
			for o := range cmd.Output() {
				c <- o
			}
			// out of the loop, the output is closed and done for this specific
			// provisioner.
			w.Done()
		}(c, cmd, w)
	}

	go func(c chan CommandOutput, w *sync.WaitGroup) {
		w.Wait()
		close(c)
	}(c, w)

	return c
}

func (cmds Commands) Wait() error {
	var errs []error
	// all of the commands must be complete, so the
	// order in which we collect the errors does not matter.
	//
	// Ie, it's not inefficient to wait for them in slice order,
	// even though they'll be completing in random order.
	for _, p := range cmds {
		if err := p.Wait(); err != nil {
			errs = append(errs, CommandError{
				CommandName: p.CommandName(),
				MachineName: p.MachineName(),
				Err:         err,
			})
		}
	}

	return errors.JoinSep(errs, "\n")
}
