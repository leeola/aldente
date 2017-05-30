package aldente

import (
	"fmt"
	"sync"

	"github.com/leeola/errors"
)

type Provisioners []Provisioner

type ProvisionersError struct {
	MachineName  string
	ProviderName string
	Err          error
}

func (ps Provisioners) Output() <-chan ProvisionOutput {
	c := make(chan ProvisionOutput, 10)
	var w sync.WaitGroup
	w.Add(len(ps))

	for _, p := range ps {
		go func(c chan ProvisionOutput, p Provisioner, w sync.WaitGroup) {
			for o := range p.Output() {
				c <- o
			}
			// out of the loop, the output is closed and done for this specific
			// provisioner.
			w.Done()
		}(c, p, w)
	}

	go func(c chan ProvisionOutput, w sync.WaitGroup) {
		w.Wait()
		close(c)
	}(c, w)

	return c
}

func (ps Provisioners) Wait() error {
	var errs []error
	// all of the provisioners must be complete, so the
	// order in which we collect the errors does not matter.
	//
	// Ie, it's not inefficient to wait for them in slice order,
	// even though they'll be completing in random order.
	for _, p := range ps {
		if err := p.Wait(); err != nil {
			errs = append(errs, ProvisionersError{
				MachineName:  p.MachineName(),
				ProviderName: p.ProviderName(),
				Err:          err,
			})
		}
	}

	return errors.JoinSep(errs, "\n")
}

func (e ProvisionersError) Error() string {
	return fmt.Sprintf("%s-%s: %s", e.ProviderName, e.MachineName, e.Err.Error())
}
