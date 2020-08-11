package main

import (
	"github.com/capeprivacy/cape/models"
	"github.com/urfave/cli/v2"
)

func init() {
	meCmd := &Command{
		Usage: "Tells you what your role is.",
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

	setCmd := &Command{
		Usage:     "Set a users role in cape.",
		Arguments: []*Argument{UserEmailArg, RoleArg},
		Examples: []*Example{
			{
				Example: "cape roles set friend@cape.com admin",
				Description: `Changes the user with email friend@cape.com to an admin
			  Possible values are admin, user`,
			},
			{
				Example: "cape roles set --project my-project friend@cape.com project-reviewer",
				Description: `Changes the user with email friend@cape.com to a project-review within my-project.
			  Possible values are project-owner, project-contributor, project-member`,
			},
		},
		Command: &cli.Command{
			Name:   "set",
			Action: handleSessionOverrides(rolesSetCmd),
			Flags: []cli.Flag{
				clusterFlag(),
				projectLabelFlag(),
			},
		},
	}

	rolesCmd := &Command{
		Usage: "Commands for querying information about roles and modifying them.",
		Command: &cli.Command{
			Name: "roles",
			Subcommands: []*cli.Command{
				setCmd.Package(),
				meCmd.Package(),
			},
		},
	}

	commands = append(commands, rolesCmd.Package())
}

func rolesSetCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	userEmail := Arguments(c.Context, UserEmailArg).(models.Email)
	roleLabel := Arguments(c.Context, RoleArg).(models.Label)

	project := c.String("project")
	if project == "" {
		err = client.SetOrgRole(c.Context, userEmail, roleLabel)
	} else {
		err = client.SetProjectRole(c.Context, userEmail, models.Label(project), roleLabel)
	}

	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	args := struct {
		Role string
		User string
	}{
		Role: roleLabel.String(),
		User: userEmail.String(),
	}

	return u.Template("Updated role to {{ .Role | bold }} for user {{ .User | bold }}\n", args)
}

func rolesMeCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	project := c.String("project")
	var role *models.Role
	if project == "" {
		role, err = client.MyRole(c.Context)
	} else {
		role, err = client.MyProjectRole(c.Context, models.Label(project))
	}

	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("Role: {{ . | bold }}\n", role.Label.String())
}
