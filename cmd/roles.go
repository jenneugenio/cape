package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	createCmd := &Command{
		Usage:     "Create a new role",
		Arguments: []*Argument{RoleLabelArg},
		Examples: []*Example{
			{
				Example:     "cape roles create data-scientist",
				Description: "Creates a new role with the label 'data-scientist'.",
			},
			{
				Example:     "CAPE_CLUSTER=prod cape roles create data-scientist",
				Description: "Creates a role on the cape cluster called prod.",
			},
		},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(rolesCreateCmd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	removeCmd := &Command{
		Usage:     "Remove command removes a role",
		Arguments: []*Argument{RoleLabelArg},
		Examples: []*Example{
			{
				Example:     "cape roles remove data-scientist",
				Description: "Removes a new role with the label 'data-scientist'.",
			},
			{
				Example:     "cape roles remove --yes data-scientist",
				Description: "Removes a role skipping the confirm dialog.",
			},
		},
		Command: &cli.Command{
			Name:   "remove",
			Action: handleSessionOverrides(rolesRemoveCmd),
			Flags: []cli.Flag{
				clusterFlag(),
				yesFlag(),
			},
		},
	}

	listCmd := &Command{
		Usage: "Lists all the roles on the cluster",
		Examples: []*Example{
			{
				Example:     "cape roles list",
				Description: "Lists all roles",
			},
		},
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(rolesListCmd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	membersCmd := &Command{
		Usage:     "Lists all the identities assigned a role",
		Arguments: []*Argument{RoleLabelArg},
		Examples: []*Example{
			{
				Example:     "cape roles members admin",
				Description: "Lists all the identities assigned role admin",
			},
		},
		Command: &cli.Command{
			Name:   "members",
			Action: handleSessionOverrides(rolesMembersCmd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	rolesCmd := &Command{
		Usage: "Commands for querying information about roles and modifying them",
		Command: &cli.Command{
			Name: "roles",
			Subcommands: []*cli.Command{
				createCmd.Package(),
				removeCmd.Package(),
				listCmd.Package(),
				membersCmd.Package(),
			},
		},
	}

	commands = append(commands, rolesCmd.Package())
}

func rolesCreateCmd(c *cli.Context) error {
	args := Arguments(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	label := args["label"].(primitives.Label)

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	// TODO specify identities
	_, err = client.CreateRole(c.Context, label, nil)
	if err != nil {
		return err
	}

	fmt.Printf("Created the role '%s'.", label)
	return nil
}

func rolesRemoveCmd(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	ui := UI(c.Context)

	args := Arguments(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	label := args["label"].(primitives.Label)
	if !skipConfirm {
		err := ui.Confirm(fmt.Sprintf("Do you really want to delete the role %s and all of its assignments?", label))
		if err != nil {
			return err
		}
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	role, err := client.GetRoleByLabel(c.Context, label)
	if err != nil {
		return err
	}

	err = client.DeleteRole(c.Context, role.ID)
	if err != nil {
		return err
	}

	fmt.Printf("The role '%s' has been deleted.", label)

	return nil
}

func rolesListCmd(c *cli.Context) error {
	ui := UI(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	roles, err := client.ListRoles(c.Context)
	if err != nil {
		return err
	}

	header := []string{"Label"}
	body := make([][]string, len(roles))
	for i, r := range roles {
		body[i] = []string{r.Label.String()}
	}

	return ui.Table(header, body)
}

func rolesMembersCmd(c *cli.Context) error {
	ui := UI(c.Context)
	args := Arguments(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	label := args["label"].(primitives.Label)

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	role, err := client.GetRoleByLabel(c.Context, label)
	if err != nil {
		return err
	}

	identities, err := client.GetMembersRole(c.Context, role.ID)
	if err != nil {
		return err
	}

	header := []string{"Type", "Email"}
	body := make([][]string, len(identities))
	for i, identity := range identities {
		typ, err := identity.GetID().Type()
		if err != nil {
			return err
		}

		typeStr := ""
		if typ == primitives.UserType {
			typeStr = primitives.UserType.String()
		} else if typ == primitives.ServicePrimitiveType {
			typeStr = primitives.ServicePrimitiveType.String()
		}

		body[i] = []string{typeStr, identity.GetEmail().String()}
	}

	return ui.Table(header, body)
}
