package autoload

import (
	ald "github.com/leeola/motley"
	"github.com/leeola/motley/autoload"
	provider "github.com/leeola/motley/providers/manual"
	cu "github.com/leeola/motley/util/configunmarshaller"
)

func init() {
	autoload.RegisterLoader(autoLoader)
}

func autoLoader(cu cu.ConfigUnmarshaller, a *ald.Aldente) error {
	ps, err := provider.FromConfigUnmarshaller(cu)
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
