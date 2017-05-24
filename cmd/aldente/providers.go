package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/leeola/aldente/autoload"
	"github.com/urfave/cli"
)

func ProvidersCmd(ctx *cli.Context) error {
	configPaths := ctx.GlobalStringSlice("config")

	if len(configPaths) <= 0 {
		return errors.New("error: at least one aldente config is required")
	}

	a, err := autoload.LoadAldente(configPaths)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\tNAME\tTYPE")

	for i, p := range a.Providers() {
		fmt.Fprintln(w, fmt.Sprintf("%d\t%s\t%s", i+1, p.Name(), p.Type()))
	}

	return w.Flush()
}
