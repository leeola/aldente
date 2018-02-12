package embed

import (
	"errors"
	"fmt"

	"github.com/leeola/motley"
)

type Config struct {
	DB         motley.Database
	Connectors []motley.Connector
	Providers  []motley.Provider

	ArchitectureConfigs []motley.ArchitectureConfig

	MachineConfigs []motley.MachineConfig
}

type Motley struct {
	config Config
	db     motley.Database

	providers  map[string]motley.Provider
	connectors map[string]motley.Connector

	archConfigs    map[string]motley.ArchitectureConfig
	machineConfigs map[string]motley.MachineConfig
}

func New(c Config) (*Motley, error) {
	if c.DB == nil {
		return nil, errors.New("missing required config: DB")
	}

	providersMap := map[string]motley.Provider{}
	for _, p := range c.Providers {
		n := p.Name()

		if _, exists := providersMap[n]; exists {
			return nil, errors.New("duplicate provider name configured")
		}

		providersMap[n] = p
	}

	connectorsMap := map[string]motley.Connector{}
	for _, p := range c.Connectors {
		n := p.Name()

		if _, exists := providersMap[n]; exists {
			return nil, errors.New("duplicate connector name configured")
		}

		connectorsMap[n] = p
	}

	archConfigs := map[string]motley.ArchitectureConfig{}
	for _, ac := range c.ArchitectureConfigs {
		if _, exists := archConfigs[ac.Machine]; exists {
			return nil, errors.New("duplicate architecture name configured")
		}

		archConfigs[ac.Machine] = ac
	}

	machineConfigs := map[string]motley.MachineConfig{}
	for _, mc := range c.MachineConfigs {
		if _, exists := machineConfigs[mc.Name]; exists {
			return nil, errors.New("duplicate machine name configured")
		}

		machineConfigs[mc.Name] = mc
	}

	mot := &Motley{
		config:         c,
		db:             c.DB,
		providers:      providersMap,
		connectors:     connectorsMap,
		archConfigs:    archConfigs,
		machineConfigs: machineConfigs,
	}

	if err := mot.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize: %s", err)
	}

	return mot, nil
}

func (m *Motley) init() error {
	// TODO(leeola): upsert all of the machine configs
	return nil
}

func (m *Motley) Status(groupName string) (motley.Status, error) {
	return motley.Status{}, nil
}
