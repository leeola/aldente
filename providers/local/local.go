package local

import (
	"encoding/json"

	ald "github.com/leeola/aldente"
	"github.com/leeola/errors"
)

const ProviderType = "local"

type Config struct {
	Name    string
	Workdir string
}

// Local implements Provider, Provisioner, and Machine for the local system.
//
// Local implements all of these because the local system has no provisioning
// steps currently. If any are added in the future, they will likely require
// Local to  split functionality into specific implementations.
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

func (l *Local) Machine(pr ald.ProviderRecord) (ald.Machine, error) {
	return nil, errors.New("not implemented")
}

func (l *Local) Provision(machineName string) (ald.Provisioner, error) {
	// local implements the provisioner
	return l, nil
}

func (l *Local) Output() <-chan ald.ProvisionOutput {
	// buffered so our write doesn't block when we send on it.
	c := make(chan ald.ProvisionOutput, 1)
	c <- ald.ProvisionOutput{
		Name:     l.config.Name,
		Provider: ProviderType,
		ProvisionStatus: ald.ProvisionStatus{
			State: ald.Provisioned,
		},
	}
	close(c)
	return c
}

func (l *Local) Wait() (ald.ProvisionHistory, ald.ProviderRecord, error) {
	h := ald.ProvisionHistory{
		ald.ProvisionStatus{
			State: ald.Provisioned,
		},
	}

	j, err := json.Marshal(l.config)
	if err != nil {
		return nil, nil, err
	}

	return h, ald.ProviderRecord(j), nil
}
