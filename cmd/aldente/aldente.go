package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/leeola/aldente"
	"github.com/leeola/aldente/databases/marshaldb"
	"github.com/urfave/cli"

	autoload "github.com/leeola/aldente/autoload"
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "registered",
					Usage: "only show registered provider implementations",
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func NewCmd(ctx *cli.Context) error {
	configPaths := ctx.GlobalStringSlice("config")
	groupName := ctx.Args().First()

	if len(configPaths) <= 0 {
		return errors.New("error: at least one aldente config is required")
	}

	if groupName == "" {
		return errors.New("error: group name is required")
	}

	// TODO(leeola): Make configurable.
	db, err := marshaldb.New(".aldente.db")
	if err != nil {
		return err
	}

	c := aldente.Config{
		Db:          db,
		ConfigPaths: configPaths,
	}
	a, err := aldente.New(c)
	if err != nil {
		return err
	}

	if err := autoload.LoadAldente(configPaths, a); err != nil {
		return err
	}

	return a.NewGroup(groupName)
}

func ListCmd(ctx *cli.Context) error {
	configPaths := ctx.GlobalStringSlice("config")

	if len(configPaths) == 0 {
		return errors.New("error: at least one aldente config is required")
	}

	// TODO(leeola): Make the database configurable. For now it's hardcoded.
	db, err := marshaldb.New(".aldente.db")
	if err != nil {
		return err
	}

	c := aldente.Config{
		Db:          db,
		ConfigPaths: configPaths,
	}
	a, err := aldente.New(c)
	if err != nil {
		return err
	}

	ms, err := a.MachineRecords()
	if err != nil {
		return err
	}

	for _, m := range ms {
		fmt.Println(m.Name, m)
	}

	return nil
}

func NotImplementedCmd(ctx *cli.Context) error {
	return errors.New("not implemented")
}
