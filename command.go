package aldente

import (
	"fmt"

	"github.com/leeola/errors"
)

type Command interface {
	// Name returns the name the running command.
	Name() string

	// MachineName returns the machine name that is running the command.
	MachineName() string

	// Start the command's execution.
	Start() error

	// Wait for the entire Command process to be done.
	Wait() error
}

// Commands convenience methods like Start and Wait for underlying Commands.
type Commands []Command

// CommandsError provides metadata about the command that resulted in an error.
type CommandsError struct {
	Name        string
	MachineName string
	Err         error
}

// Start all commands.
func (cmds Commands) Start() error {
	for _, c := range cmds {
		if err := c.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Wait for all commands.
func (cmds Commands) Wait() error {
	var errs []error
	// all of the commands must be complete, so the
	// order in which we collect the errors does not matter.
	//
	// Ie, it's not inefficient to wait for them in slice order,
	// even though they'll be completing in random order.
	for _, c := range cmds {
		if err := c.Wait(); err != nil {
			errs = append(errs, CommandsError{
				Name:        c.Name(),
				MachineName: c.MachineName(),
				Err:         err,
			})
		}
	}

	return errors.JoinSep(errs, "\n")
}

func (e CommandsError) Error() string {
	return fmt.Sprintf("%s-%s: %s", e.Name, e.MachineName, e.Err.Error())
}
