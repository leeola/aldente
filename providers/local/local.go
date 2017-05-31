package local

import (
	ald "github.com/leeola/aldente"
	"github.com/leeola/errors"
)

const ProviderType = "local"

type Config struct {
	Name    string
	Workdir string
}

type Local struct {
	config Config
}

func New(c Config) (*Local, error) {
	return &Local{
		config: c,
	}, nil
}

func (l *Local) Name() string {
	return l.config.Name
}

func (l *Local) Type() string {
	return ProviderType
}

func (l *Local) Machine(pr ald.MachineRecord) (ald.Machine, error) {
	return nil, errors.New("not implemented")
}

func (l *Local) Provision(machineName string) (ald.Provisioner, error) {
	return NewProvision(machineName, l.config)
}
