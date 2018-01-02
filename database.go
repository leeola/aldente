package motley

import (
	"encoding/json"
)

// Database implements a basic backend for Motley to write state to.
//
// Note that the database is not intended to be used from multiple
// Motley processes at the same time.
//
// The database has been designed to be a simple, dumb storage mechanism
// for Motley record types.
type Database interface {
	// // CreateMachine stores the given MachineRecord.
	// //
	// // This will fail if a machine already exists in the database with
	// // the same group name, machine name, and replication index.
	// CreateMachine(MachineRecord) error

	// // Groups lists the created groups in the database.
	// Groups() ([]string, error)

	// // GroupMachines lists the machines for the specified group.
	// GroupMachines(group string) ([]MachineRecord, error)

	// Not implemented
	// DeleteGroup(group string) error

	// // UpdateMachine updates the given machine record.
	// //
	// // Note that the group name, machine name and connection  must match an
	// // existing record or this method will return an error.
	// //
	// // Machines can only be implemented for the initial configuration when
	// // the group was first created.
	// UpdateMachine(MachineRecord) error

	// UpsertMachine creates or ignores the machine depending on if it exists.
	//
	// The provided Machine must have either Provider or Connection specified
	// in the configuration.
	UpsertMachine(MachineRecord) (upserted bool, err error)

	//UpdateTaskStates(Task)
}

// MachineRecord contains information to be able to store and reconstruct a Machine.
type MachineRecord struct {
	MachineConfig

	// Provider name of the provider that handles the machine.
	Provider string `json:"provider"`

	// ConnectionRecord is connection specific data for the given record.
	//
	// This data is used to reconstruct machines between Motley sessions. Eg,
	// it may record an ip, port, key, etc to connect to. The data stored
	// depends on the connection and/or provider.
	ConnectionRecord json.RawMessage `json:"connectionRecord"`
}
