package motley

type Connector interface {
	Name() string
	Machine() (Machine, error)
}

type Machine interface {
	Status() error
	Exec(machineName string) chan ExecOutput
	Close() error
}

type ExecOutput struct {
	Line  string
	Error error
}
