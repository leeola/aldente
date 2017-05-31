package aldente

import "github.com/leeola/errors"

// TODO(leeola): Pretty much all of the commands in the aldente package spec should
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

// Command executes the given command for the given machine.
func (a *Aldente) Command(group, commandName string) (Commands, error) {
	machineRecords, err := a.db.GroupMachines(group)
	if err != nil {
		return nil, err
	}

	commandConfig, ok := a.commandConfigs[commandName]
	if !ok {
		return nil, errors.Errorf("command not found in config: %s", commandName)
	}

	totalMachines := len(commandConfig.Machines)
	if totalMachines == 0 {
		return nil, errors.New("no machines configured for command")
	}

	commands := make(Commands, totalMachines)
	for i, mr := range machineRecords {
		p, ok := a.providers[mr.Provider]
		if !ok {
			return nil, errors.Errorf("recorded machine provider not configured: %s",
				mr.Provider)
		}

		m, err := p.Machine(mr)
		if err != nil {
			return nil, err
		}

		c, err := m.Run(commandConfig.Script)
		if err != nil {
			return nil, err
		}
		commands[i] = c
	}

	return commands, nil
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

func (a *Aldente) provisionMachine(mr MachineRecord) (Provisioner, error) {
	provider, ok := a.providers[mr.Provider]
	if !ok {
		return nil, errors.Errorf("implementation not found for provider: %s", mr.Provider)
	}

	provisioner, err := provider.Provision(mr.Name)
	if err != nil {
		return nil, err
	}

	return NewDbProvisioner(a.db, mr, provisioner), nil
}

// Provision machine(s) for the given group.
func (a *Aldente) Provision(group string) (Provisioners, error) {
	mrs, err := a.db.GroupMachines(group)
	if err != nil {
		return nil, err
	}

	ps := make([]Provisioner, len(mrs))
	for i, mr := range mrs {
		p, err := a.provisionMachine(mr)
		if err != nil {
			return nil, err
		}
		ps[i] = p
	}

	return ps, nil
}
