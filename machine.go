package aldente

// Machine that can run commands via the underlying transport.
type Machine interface {
	Name() string
	ProviderName() string
	Run(script string) (Command, error)
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
