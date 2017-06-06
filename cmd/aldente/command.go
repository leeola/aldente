package main

import (
	"errors"
	"os"

	"github.com/leeola/aldente/autoload"
	"github.com/urfave/cli"
)

func CommandCmd(ctx *cli.Context) error {
	a, err := autoload.LoadAldente(ctx.GlobalStringSlice("config"))
	if err != nil {
		return err
	}

	if len(ctx.Args()) < 2 {
		return errors.New("error: group and command required")
	}

	group := ctx.Args().Get(0)
	command := ctx.Args().Get(1)

	if group == "" {
		return errors.New("error: missing group name")
	}

	if command == "" {
		return errors.New("error: missing command name")
	}

	commands, err := a.Command(os.Stdout, group, command)
	if err != nil {
		return err
	}

	if err := commands.Start(); err != nil {
		return err
	}

	return commands.Wait()
}
