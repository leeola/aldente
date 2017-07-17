package local

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

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

func (p *Provider) Command(w io.Writer, r ald.MachineRecord, c ald.CommandConfig) error {
	// TODO(leeola): Convert to use the eventual r.History.Last() state.
	if len(r.ProviderRecord) == 0 {
		return errors.New("machine not provisioned")
	}

	// the stored providerconfig *should* be the same as the current
	// ProviderConfig, but if the user changed a value in it then we want
	// the machine to always use the same settings. So, we unmarshal what
	// was stored.
	var pConfig ProviderConfig
	if err := json.Unmarshal(r.ProviderRecord, &pConfig); err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", c.Script)
	cmd.Dir = pConfig.Workdir
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}

func (p *Provider) Machine(r ald.MachineRecord) (ald.Machine, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) Provision(w io.Writer, machineName string) (ald.ProviderRecord, error) {
	j, err := json.Marshal(p.config)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(w, "local machine provisioned")

	return ald.ProviderRecord(j), nil

}
