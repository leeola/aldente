package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/leeola/aldente"
	autoload "github.com/leeola/aldente/autoload"
	"github.com/leeola/aldente/databases/marshaldb"
	_ "github.com/leeola/aldente/providers/dockermachine/autoload"
	"github.com/urfave/cli"
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
			Name:  "ls",
			Usage: "list machines in the config",
		},
		{
			Name:   "new",
			Usage:  "create a new machine stack",
			Flags:  []cli.Flag{},
			Action: NewCmd,
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

	db, err := marshaldb.New(".aldente.db")
	if err != nil {
		return err
	}

	c := aldente.Config{
		// TODO(leeola): Make configurable.
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

	if err := a.New(groupName); err != nil {
		return err
	}

	return nil
}
