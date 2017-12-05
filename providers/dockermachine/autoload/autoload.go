package autoload

import (
	"github.com/leeola/motley"
	"github.com/leeola/motley/autoload"
	"github.com/leeola/motley/providers/dockermachine"
)

func init() {
	autoload.RegisterLoader(autoLoader)
}

func autoLoader(configPaths []string, a *motley.Aldente) error {
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
