package aldente

import "sync"

type Machine interface {
	Name() string
	Provider() Provider
	Command(name string) <-chan CommandOutput
}

type CommandOutput struct {
	CommandName  string
	MachineName  string
	ProviderName string
	Output       string
	Error        error
}

func fanCommandOutputs(cs ...<-chan CommandOutput) <-chan CommandOutput {
	f := make(chan CommandOutput, 10)
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Add(1)
		go func(c <-chan CommandOutput) {
			for o := range c {
				f <- o
			}
			wg.Done()
		}(c)
	}
	go func(wg sync.WaitGroup) {
		wg.Wait()
		close(f)
	}(wg)
	return f
}
