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

	group := ctx.Args().Get(0)
	command := ctx.Args().Get(1)

	if group == "" {
		return errors.New("missing group name")
	}

	if command == "" {
		return errors.New("missing command name")
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
