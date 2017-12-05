package autoload

import (
	ald "github.com/leeola/motley"

	"github.com/leeola/motley/autoload/registry"
	_ "github.com/leeola/motley/databases/marshaldb/regloader"
	_ "github.com/leeola/motley/providers/local/regloader"
)

func LoadAldente(configPaths []string) (*ald.Aldente, error) {
	return registry.LoadAldente(configPaths)
}
