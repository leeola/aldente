package aldente

import (
	"io"

	"github.com/leeola/errors"
)

// TODO(leeola): Pretty much all of the commands in the interface spec should contain
// context and/or channel(s) to cancel long running operations. They're being omitted
// for simplicity during prototyping/PoC.

// Provider is responsible for creating machines.
//
// This may be on the cloud, local vmware, docker, etc.
type Provider interface {
	// Name returns the configurable Name for this provider.
	//
	// Eg, you could have three providers for the type AWS which create different
	// machines; large, etc.
	Name() string

	// Type returns the constant type for this Provider.
	//
	// If Name() returns the dynamic name such as large-aws, Type() returns the
	// implementor key, such as `"aws"`.
	Type() string

	// NewMachine allocates a new machine for the give provider.
	//
	// Configuration is done via the toml config. Eg, if you want a large aws
	// instance the aws provider will be configured to use a large instance, and
	// the name of the provider will reflect that it creates a large aws instance.
	NewMachine(string) (Machine, error)
}

// Machine that can run commands via the underlying transport.
type Machine interface {
	io.Closer

	Name() string
	Provider() string
	Run(io.Reader) (io.Reader, error)
}

// MachineGroup is a collection of machines, as described in the config.
type MachineGroup map[string]Machine

type MachineGroups map[string]MachineGroup

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

type Provision interface {
}

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
		config:    c,
		db:        c.Db,
		providers: providersMap,
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
