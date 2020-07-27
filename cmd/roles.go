package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	//listCmd := &Command{
	//	Usage: "Lists all the roles on the cluster",
	//	Examples: []*Example{
	//		{
	//			Example:     "cape roles list",
	//			Description: "Lists all roles",
	//		},
	//	},
	//	Command: &cli.Command{
	//		Name:   "list",
	//		Action: handleSessionOverrides(rolesListCmd),
	//		Flags: []cli.Flag{
	//			clusterFlag(),
	//		},
	//	},
	//}

	meCmd := &Command{
		Usage: "Tells you what your role is",
		Examples: []*Example{
			{
				Example:     "cape roles me",
				Description: "Tells you what your role in the org is (admin or user)",
			},
			{
				Example:     "cape roles me --project my-project",
				Description: "Tells you what your role in the specified project is",
			},
		},
		Command: &cli.Command{
			Name:   "me",
			Action: handleSessionOverrides(rolesMeCmd),
			Flags: []cli.Flag{
				clusterFlag(),
				projectLabelFlag(),
			},
		},
	}

	rolesCmd := &Command{
		Usage: "Commands for querying information about roles and modifying them",
		Command: &cli.Command{
			Name: "roles",
			Subcommands: []*cli.Command{
				meCmd.Package(),
			},
		},
	}

	commands = append(commands, rolesCmd.Package())
}

func rolesMeCmd(c *cli.Context) error {
	project := c.String("project")

	fmt.Println("Getting role for ... ", project)

	return nil
}
