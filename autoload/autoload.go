package autoload

import (
	ald "github.com/leeola/aldente"
	"github.com/leeola/aldente/autoload/registry"
)

func LoadAldente(configPaths []string) (*ald.Aldente, error) {
	return registry.LoadAldente(configPaths)
}
