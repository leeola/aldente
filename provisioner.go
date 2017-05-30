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
	// MachineName returns the machine name that is being provisioned.
	MachineName() string

	// ProviderName returns the provider name that is provisioning the machine.
	ProviderName() string

	// Output returns a channel to monitor the progress of a provisioner.
	//
	// If an error is encountered, it can be found Wait() return value.
	Output() <-chan ProvisionOutput

	// Record returns the the ProviderRecord, blocking until it's available.
	//
	// The provisioning is not guaranteed to be finished when the Record becomes
	// available. For that, use Wait().
	Record() (ProviderRecord, error)

	// Wait for the entire Provisioner process to be done.
	//
	// If a timeout is needed to prevent waiting too long, use a timeout on the
	// output(s), while waiting for the Provisioned state.
	Wait() error
}

// ProvisionOutput contains a state and message sent during the provisioning.
type ProvisionOutput struct {
	// State is state associated with this Output.
	State ProvisionState `json:"state"`

	// Message is an optional message for the machine's state.
	Message string `json:"message,omitempty"`
}
