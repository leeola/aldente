package local

import (
	"io"
	"os/exec"

	ald "github.com/leeola/aldente"
)

type Machine struct {
	MachineRecord  ald.MachineRecord
	ProviderRecord ProviderConfig
}

func (m Machine) Name() string {
	return m.MachineRecord.Name
}

func (m Machine) GroupName() string {
	return m.MachineRecord.Group
}

func (m Machine) ProviderName() string {
	return m.MachineRecord.Provider
}

func (m Machine) Command(w io.Writer, cmdConf ald.CommandConfig) (ald.Command, error) {
	C := exec.Command("bash", "-c", cmdConf.Script)
	C.Dir = m.ProviderRecord.Workdir
	c := Command{
		name:        cmdConf.Name,
		machineName: m.Name(),
		C:           C,
	}
	c.C.Stdout = w
	c.C.Stderr = w
	return c, nil
}
