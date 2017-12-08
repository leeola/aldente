package main

import (
	"errors"

	"github.com/urfave/cli"
)

func CommandCmd(ctx *cli.Context) error {
	// a, err := autoload.LoadAldente(ctx.GlobalStringSlice("config"))
	// if err != nil {
	// 	return err
	// }

	// if len(ctx.Args()) < 2 {
	// 	return errors.New("error: group and command required")
	// }

	// group := ctx.Args().Get(0)
	// command := ctx.Args().Get(1)

	// if group == "" {
	// 	return errors.New("error: missing group name")
	// }

	// if command == "" {
	// 	return errors.New("error: missing command name")
	// }

	// return a.Command(os.Stdout, group, command)
	return errors.New("not implemented")
}
