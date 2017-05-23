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

type ProviderLoader func(cu.ConfigUnmarshaller) (ald.Provider, error)
type DatabaseLoader func(cu.ConfigUnmarshaller) (ald.Database, error)

func RegisterProvider(l ProviderLoader) error {
	if providerLoaders == nil {
		return errors.New("providers already loaded")
	}

	providerLoaders = append(providerLoaders, l)

	return nil
}

func RegisterDatabase(l DatabaseLoader) error {
	if databaseLoaders == nil {
		return errors.New("databases already loaded")
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
	return nil, nil
}

func LoadProviders(cu cu.ConfigUnmarshaller) (ald.Providers, error) {
	return nil, nil
}
