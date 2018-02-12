package motley

type Motley interface {
	Status(groupName string) (Status, error)
}

type Status struct {
	GroupName string

	Machines []Machine
}

type MachineStatus struct {
	MachineName string
}
