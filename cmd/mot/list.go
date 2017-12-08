package main

import (
	"errors"

	"github.com/urfave/cli"
)

func ListCmd(ctx *cli.Context) error {
	// configPaths := ctx.GlobalStringSlice("config")

	// if len(configPaths) <= 0 {
	// 	return errors.New("error: at least one motley config is required")
	// }

	// a, err := autoload.LoadAldente(configPaths)
	// if err != nil {
	// 	return err
	// }

	// groups, err := a.Groups()
	// if err != nil {
	// 	return err
	// }

	// return listMachines(a, groups)
	return errors.New("not implemented")
}

// func listMachines(a *ald.Aldente, groups []string) error {
// 	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
// 	fmt.Fprintln(w, "\tGROUP\tNAME\tPROVIDER")
//
// 	sort.Strings(groups)
//
// 	var i int
// 	for _, group := range groups {
// 		machines, err := a.GroupMachines(group)
// 		if err != nil {
// 			return err
// 		}
//
// 		slice.Sort(machines[:], func(i, j int) bool {
// 			return machines[i].Name < machines[j].Name
// 		})
//
// 		for _, m := range machines {
// 			i++
//
// 			// TODO(leeola): unmarshal providerrecord field into a commonfields
// 			// struct, so that we can list IP, etc, if it's available.
//
// 			fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%s\t%s", i, group, m.Name, m.Provider))
//
// 			// this is a style choice, to offer some visual grouping between each
// 			// machine group set. Only the first machine of a group will show the group
// 			// it belongs to.
// 			group = ""
// 		}
// 	}
//
// 	return w.Flush()
// }
