package aldente

import (
	"fmt"
	"sync"
)

type MultiProvision struct {
	Ps []Provisioner
}

type MultiError struct {
	MachineName  string
	ProviderName string
	Err          error
}

type MultiOutput struct {
	Output
	MachineName  string
	ProviderName string
}

func (m *MultiProvision) Output() <-chan MultiOutput {
	c := make(chan MultiOutput, 10)
	var w sync.WaitGroup
	w.Add(len(m.Ps))

	for _, p := range m.Ps {
		go func(c chan MultiOutput, p Provisioner, w sync.WaitGroup) {
			mn, pn := p.MachineName(), p.ProviderName()
			for o := range p.Output() {
				c <- MultiOutput{
					Output:       o,
					MachineName:  mn,
					ProviderName: pn,
				}
			}
			// out of the loop, the output is closed and done for this specific
			// provisioner.
			w.Done()
		}(c, p, w)
	}

	go func(c chan MultiOutput, w sync.WaitGroup) {
		w.Wait()
		close(c)
	}(c, w)

	return c
}

func (m *MultiProvision) Wait() []error {
	var errs []error
	// all of the provisioners must be complete, so the
	// order in which we collect the errors does not matter.
	//
	// Ie, it's not inefficient to wait for them in slice order,
	// even though they'll be completing in random order.
	for _, p := range m.Ps {
		if err := p.Wait(); err != nil {
			errs = append(errs, MultiError{
				MachineName:  p.MachineName(),
				ProviderName: p.ProviderName(),
				Err:          err,
			})
		}
	}

	return errs
}

func (e MultiError) Error() string {
	return fmt.Sprintf("%s-%s: %s", e.ProviderName, e.MachineName, e.Err.Error())
}
