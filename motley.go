package motley

import "io"

type Motley interface {
	UseConfig(...io.Reader) (hash string, err error)
	Status(groupName string) (Status, error)
}

type Status struct {
	GroupName string

	Machines []Machine
}

type MachineStatus struct {
	MachineName string
}
