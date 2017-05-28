package aldente

// Provisioner monitors an ongoing provisioning process.
type Provisioner interface {
	// Output returns a channel to monitor the output of a provisioner.
	Output() chan<- ProvisionOutput

	// Wait blocks until the provisioning is done.
	//
	// The returned history contains all hisories in order.
	Wait() (ProvisionHistory, ProviderRecord, error)
}

// ProvisionOutput
type ProvisionOutput struct {
	Name     string
	Provider string
	ProvisionStatus
}

// MultiProvisoner implements channel fanning and multierror provisions.
type MultiProvisioner struct {
	Provisioners []Provisioner
}

// ProvisionHistory contains order provisioning statuses, with helper methods.
type ProvisionHistory []ProvisionStatus

// ProvisionState contains some information about each provision state.
type ProvisionStatus struct {
	// State is the last known state of provisioning.
	State ProvisionState `json:"state"`

	// Message is an optional message for the machine's state.
	Message string `json:"message,omitempty"`
}

func (h *ProvisionHistory) Add(s ProvisionStatus) {
	*h = append(*h, s)
}

func (h ProvisionHistory) State() ProvisionState {
	l := len(h)
	if l == 0 {
		return Unknown
	}
	return h[l].State
}

func (h ProvisionHistory) Message() string {
	l := len(h)
	if l == 0 {
		return ""
	}
	return h[l].Message
}
