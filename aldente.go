package aldente

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/leeola/errors"
)

// TODO(leeola): Pretty much all of the commands in the interface spec should contain
// context and/or channel(s) to cancel long running operations. They're being omitted
// for simplicity during prototyping/PoC.

// Provider is responsible for creating machines.
//
// This may be on the cloud, local vmware, docker, etc.
type Provider interface {
	// Name returns the constant name for this Provider.
	Name() string

	// Type() returns the constant type for this Provider.
	Type() string

	// NewMachine allocates a new machine for the give provider.
	//
	// Configuration is done via the toml config. Eg, if you want a large aws
	// instance the aws provider will be configured to use a large instance, and
	// the name of the provider will reflect that it creates a large aws instance.
	NewMachine(string) (Machine, error)
}

// Machine is an SSH-able vm or container created by a Provider.
type Machine struct {
	Name     string
	Group    string
	Provider string
	Host     string
	SSHPort  int
}

// Resource defines a filesystem resource to be created and copied to a machine.
//
// For example a Git resource will clone the given repo to a local temp directory.
// The provision or build steps will then copy the resource into the machines
// as defined by the config.
type Resource interface {
	// Path returns the path for the given resource.
	//
	// Note that resources should lazily load, so in the case of Git it will not
	// be cloned until Path() is first called.
	Path() string
}

type Provision struct {
}

type Config struct {
	ConfigPaths []string
	Db          Database
}

type Aldente struct {
	config         Config
	db             Database
	providers      map[string]Provider
	machineConfigs []Machine
}

func New(c Config) (*Aldente, error) {
	machines, err := loadMachines(c.ConfigPaths)
	if err != nil {
		return nil, err
	}

	return &Aldente{
		config:         c,
		db:             c.Db,
		providers:      map[string]Provider{},
		machineConfigs: machines,
	}, nil
}

func (a *Aldente) Configs(name string, p Provider) error {
	return nil
}

func (a *Aldente) AddProvider(name string, p Provider) error {
	if _, exists := a.providers[name]; exists {
		return errors.Errorf("error: provider with name %q already added", name)
	}

	a.providers[name] = p
	return nil
}

// New creates a new machine from the given Provider.
func (a *Aldente) New(groupName string) error {
	groupMachines := make([]Machine, len(a.machineConfigs))
	for i, m := range a.machineConfigs {
		m.Group = groupName
		groupMachines[i] = m

		if err := a.db.Add(m); err != nil {
			return err
		}
	}

	return nil
}

// loadMachines loads the machines from the given configuration.
func loadMachines(configPaths []string) ([]Machine, error) {
	// TODO(leeola): load from multiple configs
	configPath := configPaths[0]

	f, err := os.Open(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open config")
	}
	defer f.Close()

	var conf struct {
		Machines []Machine `toml:"machines"`
	}

	if _, err := toml.DecodeReader(f, &conf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	return conf.Machines, nil
}
