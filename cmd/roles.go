package main

import (
	"github.com/capeprivacy/cape/models"
	"github.com/urfave/cli/v2"
)

func init() {
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
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	project := c.String("project")
	var role *models.Role
	if project != "" {
		role, err = client.MyRole(c.Context)
	} else {
		role, err = client.MyProjectRole(c.Context, models.Label(project))
	}

	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("Role: {{ . | bold }}!\n\n", role.Label.String())
}
