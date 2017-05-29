package main

import (
	"errors"
	"fmt"

	"github.com/leeola/aldente/autoload"
	"github.com/urfave/cli"
)

func CreateCmd(ctx *cli.Context) error {
	configPaths := ctx.GlobalStringSlice("config")
	group := ctx.Args().First()

	if len(configPaths) <= 0 {
		return errors.New("error: at least one aldente config is required")
	}

	a, err := autoload.LoadAldente(configPaths)
	if err != nil {
		return err
	}

	if err := a.CreateGroup(group); err != nil {
		return err
	}

	if !ctx.Bool("not-provision") {
		// TODO(leeola): this interface is likely to change a lot,
		// as provisioning is not yet *really* implemented, local just
		// allocates a machine interface.
		p, err := a.Provision(group)
		if err != nil {
			return err
		}

		// print the output.
		for o := range p.Output() {
			fmt.Printf(
				"[%s:%s] (%s) %s\n",
				o.Name, o.Provider,
				o.State, o.Message,
			)
		}

		fmt.Println("\ngroup provisioned machines:")
	} else {
		fmt.Println("\ngroup created machines:")
	}

	return listMachines(a, []string{group})
}
