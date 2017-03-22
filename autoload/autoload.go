package autoload

import (
	"github.com/leeola/aldente"
)

var (
	loadedConfigPaths []string
	loadedAldente     *aldente.Aldente
	loaders           []loader
)

type loader func([]string, *aldente.Aldente) error

func RegisterLoader(l loader) error {
	if loadedAldente != nil {
		if err := l(loadedConfigPaths, loadedAldente); err != nil {
			return err
		}
	}

	loaders = append(loaders, l)

	return nil
}

func LoadAldente(configPaths []string, a *aldente.Aldente) error {
	for _, l := range loaders {
		if err := l(configPaths, a); err != nil {
			return err
		}
	}

	loadedConfigPaths = configPaths
	loadedAldente = a

	return nil
}
