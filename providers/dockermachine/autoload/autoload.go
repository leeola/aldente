package autoload

import (
	"github.com/leeola/aldente"
	"github.com/leeola/aldente/autoload"
	"github.com/leeola/aldente/providers/dockermachine"
)

func init() {
	autoload.RegisterLoader(autoLoader)
}

func autoLoader(configPaths []string, a *aldente.Aldente) error {
	// TODO(leeola): load multiple configs
	ps, err := dockermachine.FromConfig(configPaths[0])
	if err != nil {
		return err
	}

	for _, p := range ps {
		if err := a.AddProvider(p.Name(), p); err != nil {
			return err
		}
	}

	return nil
}
