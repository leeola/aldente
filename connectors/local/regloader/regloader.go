package regloader

import (
	ald "github.com/leeola/motley"
	"github.com/leeola/motley/autoload/registry"
	"github.com/leeola/motley/connectors/local"
	cu "github.com/leeola/motley/util/configunmarshaller"
)

func init() {
	registry.MustRegisterConnector(Loader)
}

func Loader(cu cu.ConfigUnmarshaller) ([]ald.Connector, error) {
	var rootC struct {
		DontExpandHome   bool `toml:"dontExpandHome"`
		ConnectorConfigs []struct {
			local.Config
			// DontExpandHome *bool  `toml:"dontExpandHome"`
			Type string `toml:"type"`
		} `toml:"connection"`
	}

	if err := cu.Unmarshal(&rootC); err != nil {
		return nil, err
	}

	var cs []ald.Connector

	// create local providers for each of the configured providers
	for _, c := range rootC.ConnectorConfigs {
		if c.Type != local.ConnectorType {
			continue
		}

		// c.Workdir = util.HomeExpander(c.Workdir, rootC.DontExpandHome, c.DontExpandHome)

		conn, err := local.New(c.Config)
		if err != nil {
			return nil, err
		}

		cs = append(cs, conn)
	}

	return cs, nil
}
