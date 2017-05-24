package autoload

import (
	ald "github.com/leeola/aldente"

	"github.com/leeola/aldente/autoload/registry"
	_ "github.com/leeola/aldente/databases/marshaldb/regloader"
)

func LoadAldente(configPaths []string) (*ald.Aldente, error) {
	return registry.LoadAldente(configPaths)
}
