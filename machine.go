package aldente

import (
	"encoding/json"
	"io"
)

// Machine that can run commands via the underlying transport.
type Machine interface {
	io.Closer

	Name() string
	Provider() string
	Run(io.Reader) (io.Reader, error)
}

// MachineConfig is configuration for machines from config files.
//
// A MachineConfig is used to construct a group of machine records and
// allow providers to implement machines.
type MachineConfig struct {
	Name     string `toml:"name"`
	Provider string `toml:"provider"`
}

// MachineRecord contains information to be able to store and reconstruct a Machine.
type MachineRecord struct {
	// Name of the machine, as in the configuration.
	Name string `json:"name"`

	// Group name of the group that the machine belongs to.
	Group string `json:"group"`

	// Provider name of the provider that handles the machine.
	Provider string `json:"provider"`

	// ProvisionStatus is the last known status of provisioning.
	ProvisionStatus ProvisionStatus

	// ProviderRecord is provider specific data for the given record.
	//
	// This data is used to reconstruct machines between Aldente sessions. Eg,
	// it may record an ip, port, key, etc to connect to. The data stored
	// depends on the provider.
	ProviderRecord json.RawMessage `json:"providerRecord"`
}
