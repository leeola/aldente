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

type Provider struct {
	config Config
}

func New(c Config) (*Provider, error) {
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

func (p *Provider) Machine(pr ald.ProviderRecord) (ald.Machine, error) {
	return nil, errors.New("not implemented")
}

func (p *Provider) Provision(machineName string) (ald.Provisioner, error) {
	return nil, errors.New("not implemented")
}
