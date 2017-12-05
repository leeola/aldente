package motley

import (
	"io"

	"github.com/leeola/errors"
)

// TODO(leeola): Pretty much all of the commands in the motley package spec should
// contain context and/or channel(s) to cancel long running operations. They're
// being omitted for simplicity during prototyping/PoC.
//
// As an alternative, i may change the entire structure to be daemon based so
// that users can attach and detatch from long running deployments. Hard to say
// which direction i want to go during this prototype.

type Config struct {
	ConfigPaths    []string
	Db             Database
	CommandConfigs []CommandConfig
	MachineConfigs []MachineConfig
	Providers      []Provider
}

type Aldente struct {
	config Config
	db     Database

	// providers is a map of providers which serve to create new machines and
	// implement Machine interfaces for already existing machines.
	providers map[string]Provider

	machineConfigs []MachineConfig
	commandConfigs map[string]CommandConfig
}

func New(c Config) (*Aldente, error) {
	if c.Db == nil {
		return nil, errors.New("missing required config: Db")
	}

	if len(c.Providers) == 0 {
		return nil, errors.New("missing required config: Providers")
	}

	if len(c.MachineConfigs) == 0 {
		return nil, errors.New("missing required config: MachineConfigs")
	}

	providersMap := map[string]Provider{}
	for _, p := range c.Providers {
		n := p.Name()

		if _, exists := providersMap[n]; exists {
			return nil, errors.New("duplicate provider name configured")
		}

		providersMap[n] = p
	}

	commands := map[string]CommandConfig{}
	for _, c := range c.CommandConfigs {
		if _, exists := commands[c.Name]; exists {
			return nil, errors.New("duplicate command name in config")
		}
		commands[c.Name] = c
	}

	return &Aldente{
		config:         c,
		db:             c.Db,
		providers:      providersMap,
		machineConfigs: c.MachineConfigs,
		commandConfigs: commands,
	}, nil
}

// Groups lists groups in the db.
func (a *Aldente) Groups() ([]string, error) {
	return a.db.Groups()
}

// GroupMachines lists machines for the given group.
func (a *Aldente) GroupMachines(group string) ([]MachineRecord, error) {
	return a.db.GroupMachines(group)
}

// Commands lists the configured commands.
func (a *Aldente) Commands() []CommandConfig {
	var c []CommandConfig
	for _, v := range a.commandConfigs {
		c = append(c, v)
	}
	return c
}

func (a *Aldente) machineCommand(w io.Writer, r MachineRecord, c CommandConfig) error {
	p, ok := a.providers[r.Provider]
	if !ok {
		return errors.Errorf("recorded machine provider not configured: %s", r.Provider)
	}

	// TODO(leeola): here we will create fmtio writers and pass them to the
	// command. With each writer, we'll spin a goroutine and wait for the
	// command to be done. Once done, we will flush the fmt writer.
	//
	// This keeps the API for commands simple, while still
	// allowing more advanced write formatting like tabbed columns,
	// grouped lines, etc.
	return p.Command(w, r, c)
}

// Command executes the given command for the given machine.
func (a *Aldente) Command(w io.Writer, group, commandName string) error {
	machineRecords, err := a.db.GroupMachines(group)
	if err != nil {
		return err
	}

	commandConfig, ok := a.commandConfigs[commandName]
	if !ok {
		return errors.Errorf("command not found in config: %s", commandName)
	}

	totalMachines := len(commandConfig.Machines)
	if totalMachines == 0 {
		return errors.New("no machines configured for command")
	}

	var errs []error

	// TODO(leeola): run these concurrently.
	for _, machineName := range commandConfig.Machines {
		r, ok := findMachineRecord(machineRecords, machineName)
		if !ok {
			return errors.Errorf(
				"command configuration not found in machine record: %s", machineName)
		}

		if err := a.machineCommand(w, r, commandConfig); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.JoinSep(errs, "\n")
}

// NewGroup creates a new machine group based on the configuration.
//
// Note that the group is created, but not the actual machines. Eg, no
// VMs/Containers/etc are created from this method until
// CreateMachine(group,name) is called, only the machine records are created as
// placeholders, waiting to be created.
//
// This allows for manually allocating a machine within a group.
func (a *Aldente) CreateGroup(group string) error {
	return a.db.CreateGroup(group, a.machineConfigs)
}

// Providers lists the configured providers.
func (a *Aldente) Providers() []Provider {
	var providers []Provider
	for _, v := range a.providers {
		providers = append(providers, v)
	}
	return providers
}

func (a *Aldente) provisionMachine(w io.Writer, mRecord MachineRecord) error {
	provider, ok := a.providers[mRecord.Provider]
	if !ok {
		return errors.Errorf("implementation not found for provider: %s", mRecord.Provider)
	}

	pRecord, err := provider.Provision(w, mRecord.Name)
	if err != nil {
		return err
	}

	mRecord.ProviderRecord = pRecord
	if err := a.db.UpdateMachine(mRecord); err != nil {
		return err
	}

	return nil
}

// Provision create group and provision machine(s).
func (a *Aldente) Provision(w io.Writer, group string) error {
	if err := a.CreateGroup(group); err != nil {
		return err
	}

	mrs, err := a.db.GroupMachines(group)
	if err != nil {
		return err
	}

	var errs []error

	// TODO(leeola): run these concurrently.
	for _, mr := range mrs {
		if err := a.provisionMachine(w, mr); err != nil {
			// TODO(leeola): return a custom error object with the machine name and
			// provider included in the message and in the struct.
			errs = append(errs, err)
		}
	}

	return errors.JoinSep(errs, "\n")
}

func findMachineRecord(mrs []MachineRecord, machineName string) (MachineRecord, bool) {
	for _, mr := range mrs {
		if mr.Name == machineName {
			return mr, true
		}
	}
	return MachineRecord{}, false
}
