package regloader

import (
	ald "github.com/leeola/motley"
	"github.com/leeola/motley/autoload/registry"
	"github.com/leeola/motley/providers/local"
	"github.com/leeola/motley/util"
	cu "github.com/leeola/motley/util/configunmarshaller"
)

func init() {
	registry.MustRegisterProvider(Loader)
}

func Loader(cu cu.ConfigUnmarshaller) ([]ald.Provider, error) {
	var rootC struct {
		DontExpandHome  bool `toml:"dontExpandHome"`
		ProviderConfigs []struct {
			local.ProviderConfig
			DontExpandHome *bool  `toml:"dontExpandHome"`
			Type           string `toml:"type"`
		} `toml:"provider"`
		ProvisionConfigs []local.ProvisionConfig `toml:"provision"`
	}

	if err := cu.Unmarshal(&rootC); err != nil {
		return nil, err
	}

	var ps []ald.Provider
	// create local providers for each of the configured providers
	for _, c := range rootC.ProviderConfigs {
		if c.Type != local.ProviderType {
			continue
		}

		c.Workdir = util.HomeExpander(c.Workdir, rootC.DontExpandHome, c.DontExpandHome)

		p, err := local.New(c.ProviderConfig)
		if err != nil {
			return nil, err
		}

		ps = append(ps, p)
	}

	return ps, nil
}