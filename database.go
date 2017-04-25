package aldente

// Database implements a basic backend for Aldente to write data to.
//
// Aldente uses this to keep track of which machines were created and managed by
// aldente session.
type Database interface {
	List() ([]MachineRecord, error)
	Add(MachineRecord) error
	// Remove(MachineRecord) error
}

// MachineRecord contains information to be able to store and construct a Machine.
type MachineRecord struct {
	Name     string
	Group    string
	Provider string
	Host     string
	SSHPort  int
}
