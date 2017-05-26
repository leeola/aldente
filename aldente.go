package aldente

import "github.com/leeola/errors"

type Config struct {
	ConfigPaths    []string
	Db             Database
	MachineConfigs []MachineConfig
	Providers      []Provider
}

type Aldente struct {
	config Config
	db     Database

	// providers is a map of providers which serve to create new machines and
	// implement Machine interfaces for already existing machines.
	providers map[string]Provider

	machineConfigs []MachineConfig
}

func New(c Config) (*Aldente, error) {
	if c.Db == nil {
		return nil, errors.New("missing required config: Db")
	}

	if len(c.Providers) == 0 {
		return nil, errors.New("missing required config: Providers")
	}

	if len(c.MachineConfigs) == 0 {
		return nil, errors.New("missing required config: MachineConfigs")
	}

	providersMap := map[string]Provider{}
	for _, p := range c.Providers {
		n := p.Name()

		if _, exists := providersMap[n]; exists {
			return nil, errors.New("duplicate provider name configured")
		}

		providersMap[n] = p
	}

	return &Aldente{
		config:         c,
		db:             c.Db,
		providers:      providersMap,
		machineConfigs: c.MachineConfigs,
	}, nil
}

// Groups lists groups in the db.
func (a *Aldente) Groups() ([]string, error) {
	return a.db.Groups()
}

// GroupMachines lists machines for the given group.
func (a *Aldente) GroupMachines(group string) ([]MachineRecord, error) {
	return a.db.GroupMachines(group)
}

// NewGroup creates a new machine group based on the configuration.
//
// Note that the group is created, but not the actual machines. Eg, no
// VMs/Containers/etc are created from this method until
// CreateMachine(group,name) is called, only the machine records are created as
// placeholders, waiting to be created.
//
// This allows for manually allocating a machine within a group.
func (a *Aldente) CreateGroup(group string) error {
	return a.db.CreateGroup(group, a.machineConfigs)
}

// Providers lists the configured providers.
func (a *Aldente) Providers() []Provider {
	var providers []Provider
	for _, v := range a.providers {
		providers = append(providers, v)
	}
	return providers
}

func (a *Aldente) provisionMachine(mr MachineRecord, po ProvisionOutput) {
}

// Provision new machine(s) for the given group.
func (a *Aldente) Provision(group string) (Provision, error) {
	return nil, errors.New("not implemented")
}
