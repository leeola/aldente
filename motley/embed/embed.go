package embed

import (
	"errors"

	"github.com/leeola/motley"
)

type Config struct {
	DB         motley.Database
	Connectors []motley.Connector
	Providers  []motley.Provider
}

type Motley struct {
	config Config
	db     motley.Database

	providers  map[string]motley.Provider
	connectors map[string]motley.Connector
}

func New(c Config) (*Motley, error) {
	if c.DB == nil {
		return nil, errors.New("missing required config: DB")
	}

	if len(c.Providers) == 0 {
		return nil, errors.New("missing required config: Providers")
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
			return nil, errors.New("duplicate provider name configured")
		}

		connectorsMap[n] = p
	}

	return &Motley{
		config:     c,
		db:         c.DB,
		providers:  providersMap,
		connectors: connectorsMap,
	}, nil
}

func (m *Motley) Status(string) error {
	return nil
}
