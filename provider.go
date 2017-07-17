package aldente

import (
	"encoding/json"
	"io"
)

type ProvisionState int

const (
	Unknown ProvisionState = iota
	// Building and Built are not implemented in Aldente yet.
	//
	// // Building allows an image to be constructed for the creating state.
	// //
	// // They're positioned in front of creating/created, due to VMs/etc being
	// // based off of an image. Building constructs that image.
	// //
	// // This works for Docker style images too, and would allow a builder
	// // to construct a dockerfile based on the build instructions. The dockerfile
	// // may or may not contain step based caching, depending on the builder
	// // implementation.
	// // Building
	// // Built
	Creating
	Created
	Provisioning
	Provisioned
	Failed
)

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

	// Command runs the given commandConfig on the specified machine.
	//
	// The MachineRecord and ProviderRecord allows the provider to establish a
	// connection from the information it previously associated with the the
	// machine when it was created.
	Command(io.Writer, MachineRecord, CommandConfig) error

	// Provision based on the Provider implementation and configuration.
	//
	// Configuration is done via the toml config. Eg, if you want a large aws
	// instance the aws provider will be configured to use a large instance, and
	// the name of the provider will reflect that it creates a large aws instance.
	//
	// The actual provisioning steps, such as each ordered set of bash scripts
	// found in the `[[provisions]]` key, is again configured via the toml file
	// and upon Provider creation
	//
	// The Provider is responsible for unmarshalling all of the building,
	// provisioning etc config details needed to Provision each machine.
	Provision(w io.Writer, machineName string) (ProviderRecord, error)
}

// MachineRecord contains information to be able to store and reconstruct a Machine.
type MachineRecord struct {
	// Name of the machine, as in the configuration.
	Name string `json:"name"`

	// Group name of the group that the machine belongs to.
	Group string `json:"group"`

	// Provider name of the provider that handles the machine.
	Provider string `json:"provider"`

	// ProvisionHistory records each status in order.
	// ProvisionHistory ProvisionHistory `json:"provisionHistory"`

	// ProviderRecord is provider specific data for the given record.
	//
	// This data is used to reconstruct machines between Aldente sessions. Eg,
	// it may record an ip, port, key, etc to connect to. The data stored
	// depends on the provider.
	ProviderRecord ProviderRecord `json:"providerRecord"`
}

// ProviderRecord stores Provider data within a MachineRecord in the database.
//
// This raw json allows the Provider to instatiate Machine interfaces for a
// specific machine at any time.
type ProviderRecord json.RawMessage

func (p ProvisionState) String() string {
	switch p {
	case Unknown:
		return "Unknown"
	case Creating:
		return "Creating"
	case Created:
		return "Created"
	case Provisioning:
		return "Provisioning"
	case Provisioned:
		return "Provisioned"
	case Failed:
		return "Failed"
	default:
		return "unhandled provision state"
	}
}
