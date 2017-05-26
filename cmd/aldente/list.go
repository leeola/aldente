package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/leeola/aldente/autoload"
	"github.com/urfave/cli"
)

func ListCmd(ctx *cli.Context) error {
	configPaths := ctx.GlobalStringSlice("config")

	if len(configPaths) <= 0 {
		return errors.New("error: at least one aldente config is required")
	}

	a, err := autoload.LoadAldente(configPaths)
	if err != nil {
		return err
	}

	groups, err := a.Groups()
	if err != nil {
		return err
	}
	sort.Strings(groups)

	return listMachines(groups)
}

func listMachines(groups []string) error {
	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\tGROUP\tNAME\tPROVIDER")

	var i int
	for _, group := range groups {
		machines, err := a.GroupMachines(group)
		if err != nil {
			return err
		}

		for _, m := range machines {
			i++

			// TODO(leeola): unmarshal providerrecord field into a commonfields
			// struct, so that we can list IP, etc, if it's available.

			fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%s\t%s", i, group, m.Name, m.Provider))

			// this is a style choice, to offer some visual grouping between each
			// machine group set. Only the first machine of a group will show the group
			// it belongs to.
			group = ""
		}
	}

	return w.Flush()
}
