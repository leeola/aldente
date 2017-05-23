package aldente

import (
	"io"

	cu "github.com/leeola/aldente/util/configunmarshaller"
	"github.com/leeola/errors"
)

type Providers map[string]Provider

// func (p Providers) ProvideGroup([]MachineConfig) (MachineGroup, error) {
// 	return nil, errors.New("not implemented")
// }

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
	Providers      Providers
}

type Aldente struct {
	config Config
	db     Database

	// providers is a map of providers which serve to create new machines and
	// implement Machine interfaces for already existing machines.
	providers Providers

	machineConfigs []MachineConfig
}

func New(c Config) (*Aldente, error) {
	return &Aldente{
		config:    c,
		db:        c.Db,
		providers: c.Providers,
	}, nil
}

// loadMachineConfigs loads machine configs and checks for missing providers.
func (a *Aldente) loadMachineConfigs(cu cu.ConfigUnmarshaller) ([]MachineConfig, error) {
	var config struct {
		Machines []MachineConfig
	}
	if err := cu(&config); err != nil {
		return nil, err
	}
	ms := config.Machines

	for _, m := range ms {
		if m.Name == "" {
			return nil, errors.New("machine missing name value")
		}

		if m.Provider == "" {
			return nil, errors.Errorf("machine missing provider value: %s", m.Name)
		}

		if _, ok := a.providers[m.Provider]; !ok {
			return nil, errors.Errorf("machine's provider not implemented: %s", m.Provider)
		}
	}

	return ms, nil
}

// MachineRecords lists machines created and recorded in the db.
func (a *Aldente) MachineRecords() ([]MachineRecord, error) {
	return a.db.List()
}

// NewGroup creates a new machine group based on the configuration.
//
// Note that the group is created, but not the actual machines. Eg, no
// VMs/Containers/etc are created from this method until
// CreateMachine(group,name) is called, only the machine records are created as
// placeholders, waiting to be created.
func (a *Aldente) NewGroup(groupName string) error {
	cu := cu.New(a.config.ConfigPaths)

	machineRecords, err := a.MachineRecords()
	if err != nil {
		return err
	}

	// confirm the new group name is unique
	for _, mr := range machineRecords {
		if mr.Group == groupName {
			return errors.Errorf("group name already in use: %s", groupName)
		}
	}

	machineConfigs, err := a.loadMachineConfigs(cu)
	if err != nil {
		return err
	}

	// create a record for each machineConfig
	for _, mc := range machineConfigs {
		mr := MachineRecord{
			Name:     mc.Name,
			Group:    groupName,
			Provider: mc.Provider,
		}

		if err := a.db.Add(mr); err != nil {
			return errors.Wrapf(err,
				"failed to store record for machine config: %s", mc.Name)
		}
	}

	return nil
}
