package registry

import (
	"errors"
	"fmt"

	"github.com/leeola/motley"
	"github.com/leeola/motley/connectors/local"
	"github.com/leeola/motley/motley/embed"
	cu "github.com/leeola/motley/util/configunmarshaller"
)

var (
	// If the loaders are nil, loaders have already been produced and the
	// the slices have had their memory freed.
	providerLoaders  []ProviderLoader
	connectorLoaders []ConnectorLoader
	databaseLoaders  []DatabaseLoader
)

func init() {
	// init the loaders so they're not nil. Nil loader slice represents a freed
	// slice of loaders.
	providerLoaders = []ProviderLoader{}
	connectorLoaders = []ConnectorLoader{}
	databaseLoaders = []DatabaseLoader{}
}

type ProviderLoader func(cu.ConfigUnmarshaller) ([]motley.Provider, error)
type ConnectorLoader func(cu.ConfigUnmarshaller) ([]motley.Connector, error)
type DatabaseLoader func(cu.ConfigUnmarshaller) (motley.Database, error)

func MustRegisterConnector(l ConnectorLoader) error {
	if connectorLoaders == nil {
		panic("connectors already loaded")
	}

	connectorLoaders = append(connectorLoaders, l)

	return nil
}

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
func LoadAldente(configPaths []string) (motley.Motley, error) {
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

	conns, err := LoadConnectors(cu)
	if err != nil {
		return nil, err
	}

	var conf struct {
		Machines []motley.MachineConfig `toml:"machine"`
		Commands []motley.CommandConfig `toml:"command"`
	}
	if err := cu.Unmarshal(&conf); err != nil {
		return nil, err
	}

	eConf := embed.Config{
		// ConfigPaths:    configPaths,
		DB:         db,
		Providers:  p,
		Connectors: conns,
		// MachineConfigs: conf.Machines,
		// CommandConfigs: conf.Commands,
	}
	return embed.New(eConf)
}

func LoadDatabase(cu cu.ConfigUnmarshaller) (motley.Database, error) {
	var dbs []motley.Database
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

func LoadProviders(cu cu.ConfigUnmarshaller) ([]motley.Provider, error) {
	var ps []motley.Provider
	for _, l := range providerLoaders {
		p, err := l(cu)
		if err != nil {
			return nil, err
		}

		ps = append(ps, p...)
	}

	return ps, nil
}

func LoadConnectors(cu cu.ConfigUnmarshaller) ([]motley.Connector, error) {
	var (
		conns    []motley.Connector
		hasLocal bool
	)
	for _, l := range connectorLoaders {
		loadedConns, err := l(cu)
		if err != nil {
			return nil, err
		}

		for _, loadedConn := range loadedConns {
			if loadedConn.Name() == "local" {
				hasLocal = true
			}
		}

		conns = append(conns, loadedConns...)
	}

	// As one of the few convience "magic" features, lets always give
	// Motley a local connector, so it can be available even if
	// unconfigured.
	//
	// Useful for local builds, etc.
	if !hasLocal {
		conn, err := local.New(local.Config{Name: "local"})
		if err != nil {
			return nil, fmt.Errorf("loadconnectors: magic local connector: %s", err)
		}

		conns = append(conns, conn)
	}

	return conns, nil
}
