package aldente

// Database implements a basic backend for Aldente to write data to.
//
// Aldente uses this to keep track of which machines were created and managed by
// aldente session.
type Database interface {
	// List() (Machine, error)
	Add(Machine) error
	// Remove(Machine) error
}
