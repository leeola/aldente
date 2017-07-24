package aldente

type TaskState int

const (
	TaskStateUnknown TaskState = iota
	TaskUnstarted
	TaskInProgress
	TaskSuccess
	TaskFailure
)

type TaskType int

const (
	TaskTypeUnknown TaskType = iota
	TaskCommand
	TaskCreate
)

type Task struct {
	Id        int
	TaskType  TaskType
	TaskState TaskState
	GroupName string

	CommandTask *CommandTask `json:"commandTask"`
	// CreateTask  *CreateTask  `json:"CreateTask"`
}

type CommandTask struct {
	CommandName     string
	MachineCommands []MachineCommandTask
}

type MachineCommandTask struct {
	CommandName string
	LogKey      string
}

// type CreateTask struct {
// 	CreateMachine []CreateMachineTask
// }
//
// type CreateMachineTask struct {
// 	Build     BuildTask
// 	Provision ProvisionTask
// }
//
