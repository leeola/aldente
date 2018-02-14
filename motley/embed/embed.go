package embed

import (
	"errors"
	"fmt"
	"io"

	"github.com/leeola/motley"
)

type Config struct {
	DB         motley.Database
	Connectors []motley.Connector
}

type Motley struct {
	db motley.Database

	connectors map[string]motley.Connector
}

func New(c Config) (*Motley, error) {
	if c.DB == nil {
		return nil, errors.New("missing required config: DB")
	}

	connectorsMap := map[string]motley.Connector{}
	for _, p := range c.Connectors {
		n := p.Name()

		if _, exists := connectorsMap[n]; exists {
			return nil, errors.New("duplicate connector name configured")
		}

		connectorsMap[n] = p
	}

	mot := &Motley{
		db:         c.DB,
		connectors: connectorsMap,
	}

	if err := mot.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize: %s", err)
	}

	return mot, nil
}

// init actual machines from the given machineConfigs
func (m *Motley) init() error {
	// TODO(leeola): upsert all of the machine configs
	return nil
}

func (m *Motley) UseConfig(rs ...io.Reader) error {
	return errors.New("not implemented")
}

func (m *Motley) Status(groupName string) (motley.Status, error) {
	return motley.Status{}, nil
}
