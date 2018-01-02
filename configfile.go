package motley

// ArchitectureConfig is configuration for machines from config files.
//
// An ArchitectureConfig is used to construct a group of machine records
// and allow providers to implement machines.
type ArchitectureConfig struct {
	Machine  string `toml:"machine"`
	Provider string `toml:"provider"`
}

type MachineConfig struct {
	// Name of the machine, as in the configuration.
	Name string `toml:"name" json:"name"`

	// Group name of the group that the machine belongs to.
	Group string `toml:"group" json:"group"`

	Connection string `toml:"connection" json:"connection"`

	// Index specifies the index of this machine within the list of machines
	// matching this same name and group.
	Index int `toml:"replicationIndex" json:"replicationIndex"`
}

type CommandConfig struct {
	Name     string   `toml:"name"`
	Machines []string `toml:"machines"`
	Script   string   `toml:"script"`
}
