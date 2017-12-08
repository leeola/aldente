package motley

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

// Provider is responsible for creating machines.
//
// This may be on the cloud, local vmware, docker, etc.
type Provider interface {
	// // Create the given machine from the loaded config.
	// //
	// // The provided Group and Machine name can optionally be used by the
	// // provider to identify the created machine.
	// //
	// // The returned channel can communicate progress, as well as communicate
	// // the final value. The last value of the channel must contain the
	// // ProviderRecord or Error for the operation.
	// // See CreateOutput for further documentation on CreateOutput behavior.
	// //
	// // NOTE: The returned channel has implicit behavior because it fans into
	// // the Provision channels. Provider is a lower level implementation behind
	// // motley.
	// Create(groupName, machineName string) <-chan CreateOutput

	// Machine return an interface for the requested Machine.
	//
	// Machine(ProviderRecord) (Machine, error)

	// Name returns the configurable Name for this provider.
	//
	// Eg, you could have three providers for the type AWS which create different
	// machines; large, etc.
	Name() string

	// Type returns the constant type for this Provider.
	//
	// If Name() returns the dynamic name such as large-aws, Type() returns the
	// implementor key, such as `"aws"`.
	Type() string
}

type CreateOutput struct {
	Line           string
	ProviderRecord ProviderRecord
	Error          error
}

type ProvisionOutput struct {
	Line  string
	State ProvisionState
}

func (p ProvisionState) String() string {
	switch p {
	case Unknown:
		return "Unknown"
	case Creating:
		return "Creating"
	case Created:
		return "Created"
	case Provisioning:
		return "Provisioning"
	case Provisioned:
		return "Provisioned"
	case Failed:
		return "Failed"
	default:
		return "unhandled provision state"
	}
}
