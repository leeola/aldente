package registry

import (
	"errors"

	ald "github.com/leeola/aldente"
	cu "github.com/leeola/aldente/util/configunmarshaller"
)

var (
	// If the loaders are nil, loaders have already been produced and the
	// the slices have had their memory freed.
	providerLoaders []ProviderLoader
	databaseLoaders []DatabaseLoader
)

func init() {
	// init the loaders so they're not nil. Nil loader slice represents a freed
	// slice of loaders.
	providerLoaders = []ProviderLoader{}
	databaseLoaders = []DatabaseLoader{}
}

type ProviderLoader func(cu.ConfigUnmarshaller) ([]ald.Provider, error)
type DatabaseLoader func(cu.ConfigUnmarshaller) (ald.Database, error)

func MustRegisterProvider(l ProviderLoader) error {
	if providerLoaders == nil {
		panic("providers already loaded")
	}

	providerLoaders = append(providerLoaders, l)

	return nil
}

func MustRegisterDatabase(l DatabaseLoader) error {
	if databaseLoaders == nil {
		panic("databases already loaded")
	}

	databaseLoaders = append(databaseLoaders, l)

	return nil
}

// LoadFixity from the given configunmarshaller.
//
// Note that LoadFixity purges the registered fixities if successful.
func LoadAldente(configPaths []string) (*ald.Aldente, error) {
	if len(configPaths) == 0 {
		return nil, errors.New("a configPath is required")
	}

	cu := cu.New(configPaths)
	db, err := LoadDatabase(cu)
	if err != nil {
		return nil, err
	}

	p, err := LoadProviders(cu)
	if err != nil {
		return nil, err
	}

	aConf := ald.Config{
		ConfigPaths: configPaths,
		Db:          db,
		Providers:   p,
	}
	return ald.New(aConf)
}

func LoadDatabase(cu cu.ConfigUnmarshaller) (ald.Database, error) {
	var dbs []ald.Database
	for _, l := range databaseLoaders {
		db, err := l(cu)
		if err != nil {
			return nil, err
		}

		dbs = append(dbs, db)
	}

	if len(dbs) == 0 {
		return nil, errors.New("no database defined in configs")
	}

	if len(dbs) > 1 {
		return nil, errors.New("multiple databases defined in configs")
	}

	return dbs[0], nil
}

func LoadProviders(cu cu.ConfigUnmarshaller) ([]ald.Provider, error) {
	var ps []ald.Provider
	for _, l := range providerLoaders {
		p, err := l(cu)
		if err != nil {
			return nil, err
		}

		ps = append(ps, p...)
	}

	return ps, nil
}
