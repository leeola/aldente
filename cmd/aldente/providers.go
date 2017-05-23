package main

import (
	"errors"

	"github.com/leeola/aldente"
	"github.com/leeola/aldente/autoload"
	"github.com/leeola/aldente/databases/marshaldb"
	"github.com/urfave/cli"
)

func ProvidersCmd(ctx *cli.Context) error {
	configPaths := ctx.GlobalStringSlice("config")
	groupName := ctx.Args().First()

	if len(configPaths) <= 0 {
		return errors.New("error: at least one aldente config is required")
	}

	if groupName == "" {
		return errors.New("error: group name is required")
	}

	// TODO(leeola): Make configurable.
	db, err := marshaldb.New(".aldente.db")
	if err != nil {
		return err
	}

	c := aldente.Config{
		Db:          db,
		ConfigPaths: configPaths,
	}
	a, err := aldente.New(c)
	if err != nil {
		return err
	}

	if err := autoload.LoadAldente(configPaths, a); err != nil {
		return err
	}

	return a.NewGroup(groupName)
}
