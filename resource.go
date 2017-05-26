package aldente

// Resource defines a filesystem resource to be created and copied to a machine.
//
// For example a Git resource will clone the given repo to a local temp directory.
// The provision or build steps will then copy the resource into the machines
// as defined by the config.
type Resource interface {
	// Path returns the path for the given resource.
	//
	// Note that resources should lazily load, so in the case of Git it will not
	// be cloned until Path() is first called.
	Path() string
}
