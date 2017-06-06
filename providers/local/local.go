package local

import (
	"encoding/json"
	"fmt"
	"io"

	ald "github.com/leeola/aldente"
	"github.com/leeola/errors"
)

const ProviderType = "local"

type ProviderConfig struct {
	Name             string            `toml:"name" json:"-"`
	Workdir          string            `toml:"workdir" json:"workdir"`
	ProvisionConfigs []ProvisionConfig `toml:"-" json:"-"`
}

type ProvisionConfig struct {
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
	if err := json.Unmarshal(mr.ProviderRecord, &providerRecord); err != nil {
		return nil, err
	}

	return &Machine{
		MachineRecord:  mr,
		ProviderRecord: providerRecord,
	}, nil
}

func (p *Provider) Provision(w io.Writer, machineName string) (ald.ProviderRecord, error) {
	j, err := json.Marshal(p.config)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "local machine provisioned")

	return ald.ProviderRecord(j), nil

}
