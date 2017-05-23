package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"

	_ "github.com/leeola/aldente/providers/manual/autoload"
)

func main() {
	app := cli.NewApp()
	app.Name = "aldente"
	app.Usage = "less mushy spaghetti deployments"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "config, c",
			Usage:  "load config(s) from `PATH`",
			EnvVar: "ALDENTE_CONFIGS",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "command, c",
			Usage:  "run commands on the given group",
			Action: NotImplementedCmd,
		},
		{
			Name:   "ls",
			Usage:  "list machines in the config",
			Action: NotImplementedCmd,
		},
		{
			Name:   "new",
			Usage:  "create a new machine group",
			Flags:  []cli.Flag{},
			Action: NotImplementedCmd,
		},
		{
			Name:   "providers",
			Usage:  "list configured providers",
			Action: NotImplementedCmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func NotImplementedCmd(ctx *cli.Context) error {
	return errors.New("not implemented")
}
