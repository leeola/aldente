package autoload

import (
	"github.com/leeola/motley"

	"github.com/leeola/motley/autoload/registry"
	// _ "github.com/leeola/motley/databases/marshaldb/regloader"
	// _ "github.com/leeola/motley/providers/local/regloader"
)

func Motley(configPaths []string) (motley.Motley, error) {
	return registry.LoadAldente(configPaths)
}
