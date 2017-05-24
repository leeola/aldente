package regloader

import (
	"github.com/fatih/structs"
	ald "github.com/leeola/aldente"
	"github.com/leeola/aldente/autoload/registry"
	"github.com/leeola/aldente/databases/marshaldb"
	cu "github.com/leeola/aldente/util/configunmarshaller"
	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	registry.MustRegisterDatabase(Loader)
}

func Loader(cu cu.ConfigUnmarshaller) (ald.Database, error) {
	var c struct {
		DontExpandHome bool             `toml:"dontExpandHome"`
		Config         marshaldb.Config `toml:"marshalDatabase"`
	}

	if err := cu.Unmarshal(&c); err != nil {
		return nil, err
	}

	// if the config isn't defined, do not load anything. This is allowed.
	if structs.IsZero(c.Config) {
		return nil, nil
	}

	if !c.DontExpandHome {
		if c.Config.Path != "" {
			p, err := homedir.Expand(c.Config.Path)
			if err != nil {
				return nil, err
			}
			c.Config.Path = p
		}
	}

	return marshaldb.New(c.Config)
}
