package aldente

// MachineConfig is configuration for machines from config files.
//
// A MachineConfig is used to construct a group of machine records and
// allow providers to implement machines.
type MachineConfig struct {
	Name     string `toml:"name"`
	Provider string `toml:"provider"`
}
