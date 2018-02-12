package main

import (
	"github.com/leeola/motley/autoload"
	"github.com/urfave/cli"
)

func StatusCmd(ctx *cli.Context) error {
	groupName := ctx.Args().First()

	m, err := autoload.Motley(ctx.GlobalStringSlice("config"))
	if err != nil {
		return err
	}

	if _, err := m.Status(groupName); err != nil {
		return err
	}

	return nil
}
