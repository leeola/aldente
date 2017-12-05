package manual

import (
	ald "github.com/leeola/motley"
	cu "github.com/leeola/motley/util/configunmarshaller"
	"github.com/leeola/errors"
)

const ProviderType = "manual"

type Config struct {
	Name string
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

func (p *Provider) NewMachine(string) (ald.Machine, error) {
	return nil, errors.New("not implemented")
}

func FromConfigUnmarshaller(cu cu.ConfigUnmarshaller) ([]*Provider, error) {
	var conf struct {
		Configs []struct {
			Config
			Type string
		} `toml:"providers"`
	}

	if err := cu(&conf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	var ps []*Provider
	for _, c := range conf.Configs {
		// skip any configs that aren't for this provider
		if c.Type != ProviderType {
			continue
		}

		p, err := New(c.Config)
		if err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}

	return ps, nil
}
