package aldente

import "encoding/json"

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

	// Machine instantiates a Machine interface for the given MachineRecord.
	//
	// The ProviderRecord allows the provider to establish a connection from
	// the information it previously associated with the the machine when it was
	// created.
	Machine(ProviderRecord) (Machine, error)

	// Provision creates a new machine and then provisions it.
	//
	// Configuration is done via the toml config. Eg, if you want a large aws
	// instance the aws provider will be configured to use a large instance, and
	// the name of the provider will reflect that it creates a large aws instance.
	Provision() (Provision, error)
}

type Provision interface {
	Output() chan<- ProvisionOutput
	Wait() (ProvisionStatus, ProvisionRecord, error)
}

type ProvisionOutput struct {
	Name            string
	Provider        string
	ProvisionStatus ProvisionStatus
	Message         string
}

type ProvisionRecord json.RawMessage

// MultiProvison implements channel fanning and multierror provisions.
type MultiProvision struct {
	Provisions []Provision
}
