package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "mot"
	app.Usage = "tools for a motley architecture"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "config, c",
			Usage:  "load config(s) from `PATH`",
			EnvVar: "MOTLEY_CONFIGS",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "status",
			Usage:  "summarize status of a group of machines",
			Action: StatusCmd,
		},
		// {
		// 	Name:   "command",
		// 	Usage:  "run commands defined in the configuration",
		// 	Action: CommandCmd,
		// },
		// {
		// 	Name:   "commands",
		// 	Usage:  "list commands defined in the configuration",
		// 	Action: CommandsCmd,
		// },
		// {
		// 	Name:   "create",
		// 	Usage:  "create a new machine group",
		// 	Action: CreateCmd,
		// 	Flags: []cli.Flag{
		// 		cli.BoolFlag{
		// 			Name:  "not-provision",
		// 			Usage: "do not provision the machines in the group",
		// 		},
		// 	},
		// },
		// {
		// 	Name:   "ls",
		// 	Usage:  "list machines in the config",
		// 	Action: ListCmd,
		// },
		// {
		// 	Name:   "providers",
		// 	Usage:  "list configured providers",
		// 	Action: ProvidersCmd,
		// },
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
