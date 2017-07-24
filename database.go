package aldente

import (
	"encoding/json"
)

// Database implements a basic backend for Aldente to write data to.
//
// Aldente uses this to keep track of which machines were created and managed by
// aldente session.
type Database interface {
	// Groups lists the created groups in the database.
	Groups() ([]string, error)

	// GroupMachines lists the machines for the specified group.
	GroupMachines(group string) ([]MachineRecord, error)

	// CreateGroup and MachineRecords from the given MachineConfig slice.
	//
	// Note that this just stores the machine records without any details about
	// which ip/etc the machinerecord is for. That data will be updated via
	// UpdateMachine() once the provider is allocating the resource.
	//
	// Order of machines is not required to be consistent.
	CreateGroup(group string, machines []MachineConfig) error

	// Not implemented
	// DeleteGroup(group string) error

	// TODO(leeola): I need to flesh out the idea of updating a group's config.
	// Eg, i want to be able to add machines to a config and have them show up as
	// unallocated machines in a given group. The concern there of course, is
	// changing existing machines - as they've already been allocated.
	// Likely, updating won't be supported, and groups/machines will be strictly
	// products of their configs at the time of creation.
	//
	// On that note, storing the config records that created the machine seems
	// handy, or at the very least an identifiable git hash/ver of the config.
	// UpdateGroup(group string, []MachineConfig)

	// UpdateMachine updates the given machine record.
	//
	// Note that the Name, Group and Provider must match an existing record
	// or this method will return an error.
	//
	// Machines can only be implemented for the initial configuration when
	// the group was first created.
	UpdateMachine(MachineRecord) error

	//UpdateTaskStates(Task)
}

// MachineRecord contains information to be able to store and reconstruct a Machine.
type MachineRecord struct {
	// Name of the machine, as in the configuration.
	Name string `json:"name"`

	// Group name of the group that the machine belongs to.
	Group string `json:"group"`

	// Provider name of the provider that handles the machine.
	Provider string `json:"provider"`

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
