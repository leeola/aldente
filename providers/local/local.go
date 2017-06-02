package local

import (
	"encoding/json"

	ald "github.com/leeola/aldente"
	"github.com/leeola/errors"
)

const ProviderType = "local"

type ProviderConfig struct {
	Name             string            `toml:"name"`
	Workdir          string            `toml:"workdir"`
	ProvisionConfigs []ProvisionConfig `toml:"-"`
}

type Provider struct {
	config ProviderConfig
}

func New(c ProviderConfig) (*Provider, error) {
	return &Provider{
		config: c,
	}, nil
}

func (p *Provider) Name() string {
	return p.config.Name
}

func (p *Provider) Type() string {
	return ProviderType
}

func (p *Provider) Machine(mr ald.MachineRecord) (ald.Machine, error) {
	if len(mr.ProviderRecord) == 0 {
		return nil, errors.New("machine not provisioned")
	}

	// the stored providerconfig *should* be the same as the current
	// ProviderConfig, but if the user changed a value in it then we want
	// the machine to always use the same settings. So, we unmarshal what
	// was stored.
	var providerRecord ProviderConfig
	if err := json.Unmarshal(&providerRecord, mr.ProviderRecord); err != nil {
		return nil, err
	}

	return &Machine{
		MachineRecord:  mr,
		ProviderRecord: providerRecord,
	}, nil
}

func (p *Provider) Provision(machineName string) (ald.Provisioner, error) {
	return NewProvision(machineName, p.config, p.config.ProvisionConfigs)
}
