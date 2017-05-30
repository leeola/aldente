package aldente

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

// Provisioner monitors an ongoing provisioning process.
type Provisioner interface {
	MachineName() string

	ProviderName() string

	// Output returns a channel to monitor the progress of a provisioner.
	//
	// If an error is encountered, it can be found in the Record() channel.
	Output() <-chan Output

	// Record returns the the ProviderRecord once it's available.
	//
	// Likely, after the provisioning is entirely done.
	Record() (ProviderRecord, error)

	Wait() error
}

// Output
type Output struct {
	// State is state associated with this Output.
	State ProvisionState `json:"state"`

	// Message is an optional message for the machine's state.
	Message string `json:"message,omitempty"`
}
