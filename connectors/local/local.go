package local

import "github.com/leeola/motley"

const ConnectorType = "local"

type Config struct {
	Name string `toml:"name" json:"-"`
}

type Connector struct {
	config Config
}

func New(c Config) (*Connector, error) {
	return &Connector{
		config: c,
	}, nil
}

func (p *Connector) Name() string {
	return p.config.Name
}

func (p *Connector) Type() string {
	return ConnectorType
}

func (p *Connector) Machine() (motley.Machine, error) {
	return p, nil
}

func (p *Connector) Status() error {
	return nil
}

// Close is a noop for the local Connector
func (p *Connector) Close() error {
	return nil
}

// func (p *Connector) Command(w io.Writer, r ald.MachineRecord, c ald.CommandConfig) error {
// 	// TODO(leeola): Convert to use the eventual r.History.Last() state.
// 	if len(r.ProviderRecord) == 0 {
// 		return errors.New("machine not provisioned")
// 	}
//
// 	// the stored providerconfig *should* be the same as the current
// 	// ProviderConfig, but if the user changed a value in it then we want
// 	// the machine to always use the same settings. So, we unmarshal what
// 	// was stored.
// 	var pConfig ConnectorConfig
// 	if err := json.Unmarshal(r.ProviderRecord, &pConfig); err != nil {
// 		return err
// 	}
//
// 	cmd := exec.Command("bash", "-c", c.Script)
// 	cmd.Dir = pConfig.Workdir
// 	cmd.Stdout = w
// 	cmd.Stderr = w
// 	return cmd.Run()
// }
//
// func (p *Connector) Provision(w io.Writer, machineName string) (ald.ProviderRecord, error) {
// 	j, err := json.Marshal(p.config)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	fmt.Fprintln(w, "local machine provisioned")
//
// 	return ald.ProviderRecord(j), nil
//
// }
