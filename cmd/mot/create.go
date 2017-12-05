package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/leeola/motley/autoload"
	"github.com/urfave/cli"
)

func CreateCmd(ctx *cli.Context) error {
	configPaths := ctx.GlobalStringSlice("config")
	group := ctx.Args().First()

	if len(configPaths) <= 0 {
		return errors.New("error: at least one motley config is required")
	}

	a, err := autoload.LoadAldente(configPaths)
	if err != nil {
		return err
	}

	if !ctx.Bool("not-provision") {
		if err := a.Provision(os.Stdout, group); err != nil {
			return err
		}

		fmt.Println("\ngroup provisioned machines:")
	} else {
		if err := a.CreateGroup(group); err != nil {
			return err
		}

		fmt.Println("\ngroup created machines:")
	}

	return listMachines(a, []string{group})
}
