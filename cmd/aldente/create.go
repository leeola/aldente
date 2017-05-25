package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

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

	machines, err := a.GroupMachines(group)
	if err != nil {
		return err
	}

	fmt.Println("Group created with placeholder machines:")

	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\tNAME\tPROVIDER")

	for i, m := range machines {
		fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%s", i+1, m.Name, m.Provider))
	}

	return w.Flush()

	return nil
}
