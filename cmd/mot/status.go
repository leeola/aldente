package main

import (
	"github.com/leeola/motley/autoload"
	"github.com/urfave/cli"
)

func StatusCmd(ctx *cli.Context) error {
	_, err := autoload.LoadAldente(ctx.GlobalStringSlice("config"))
	if err != nil {
		return err
	}

	return nil
}
