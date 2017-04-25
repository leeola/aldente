package autoload

import (
	ald "github.com/leeola/aldente"
	"github.com/leeola/aldente/autoload"
	provider "github.com/leeola/aldente/providers/manual"
	cu "github.com/leeola/aldente/util/configunmarshaller"
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
