package local

import (
	"encoding/json"

	ald "github.com/leeola/aldente"
)

type ProvisionConfig struct {
}

type Provision struct {
	config           ProviderConfig
	provisionConfigs []ProvisionConfig
	machineName      string
}

func NewProvision(machineName string, pdrConf ProviderConfig, psnConf []ProvisionConfig) (*Provision, error) {
	return &Provision{
		config:           pdrConf,
		provisionConfigs: psnConf,
		machineName:      machineName,
	}, nil
}

func (l *Provision) MachineName() string {
	return l.machineName
}

func (l *Provision) ProviderName() string {
	return l.config.Name
}

func (l *Provision) Output() <-chan ald.ProvisionOutput {
	// buffered so our write doesn't block when we send on it.
	c := make(chan ald.ProvisionOutput, 1)
	c <- ald.ProvisionOutput{
		MachineName:  l.MachineName(),
		ProviderName: l.ProviderName(),
		State:        ald.Provisioned,
	}
	close(c)
	return c
}

func (l *Provision) Record() (ald.ProviderRecord, error) {
	j, err := json.Marshal(l.config)
	if err != nil {
		return nil, err
	}

	return ald.ProviderRecord(j), nil
}

func (l *Provision) Wait() error {
	// the local provisioner does no actual work.
	return nil
}
