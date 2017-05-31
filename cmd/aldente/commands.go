package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/leeola/aldente/autoload"
	"github.com/urfave/cli"
)

func CommandsCmd(ctx *cli.Context) error {
	a, err := autoload.LoadAldente(ctx.GlobalStringSlice("config"))
	if err != nil {
		return err
	}

	trimScriptTo := 25

	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\tCOMMAND\tMACHINES\tSCRIPT")
	for i, c := range a.Commands() {
		var script string
		if len(c.Script) < trimScriptTo {
			script = c.Script
		} else {
			script = c.Script[:trimScriptTo] + "..."
		}

		fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%s\t%q", i+1,
			c.Name, strings.Join(c.Machines, ","), script))
	}

	return w.Flush()
}
