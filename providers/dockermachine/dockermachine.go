package dockermachine

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/leeola/motley"
	"github.com/leeola/errors"
)

const ProviderType = "dockermachine"

type Config struct {
	Name string
	Type string
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

func (p *Provider) NewMachine(string) (motley.Machine, error) {
	return nil, errors.New("not implemented")
}

func FromConfig(configPath string) ([]*Provider, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf struct {
		Configs []Config `toml:"providers"`
	}

	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	var ps []*Provider
	for _, c := range conf.Configs {
		// skip any configs that aren't for this provider
		if c.Type != ProviderType {
			continue
		}

		p, err := New(c)
		if err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}

	return ps, nil
}
