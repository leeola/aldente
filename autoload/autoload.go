package autoload

import (
	"errors"

	ald "github.com/leeola/aldente"
	cu "github.com/leeola/aldente/util/configunmarshaller"
)

var (
	loadedConfigPaths []string
	loadedAldente     *ald.Aldente
	loaders           []loader
)

type loader func(cu.ConfigUnmarshaller, *ald.Aldente) error

func RegisterLoader(l loader) error {
	if loadedAldente != nil {
		cu := cu.New(loadedConfigPaths)

		if err := l(cu, loadedAldente); err != nil {
			return err
		}
	}

	loaders = append(loaders, l)

	return nil
}

func LoadAldente(configPaths []string, a *ald.Aldente) error {
	if loadedAldente != nil {
		return errors.New("autoload already loaded Aldente")
	}

	cu := cu.New(configPaths)

	for _, l := range loaders {
		if err := l(cu, a); err != nil {
			return err
		}
	}

	loadedConfigPaths = configPaths
	loadedAldente = a

	return nil
}
