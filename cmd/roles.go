package main

import (
	"fmt"
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/database"
	"github.com/capeprivacy/cape/primitives"
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
				membersFlag(),
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
	members := c.StringSlice("member")

	label := Arguments(c.Context, RoleLabelArg).(primitives.Label)

	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	emails := make([]primitives.Email, len(members))
	for i, emailStr := range members {
		email, err := primitives.NewEmail(emailStr)
		if err != nil {
			return err
		}
		emails[i] = email
	}

	identities, err := client.GetIdentities(c.Context, emails)
	if err != nil {
		return err
	}

	identityIDs := make([]database.ID, len(identities))
	for i, identity := range identities {
		identityIDs[i] = identity.GetID()
	}

	_, err = client.CreateRole(c.Context, label, identityIDs)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("Created the role {{ . | bold }}.\n", label.String())
}

func rolesRemoveCmd(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	label := Arguments(c.Context, RoleLabelArg).(primitives.Label)
	if !skipConfirm {
		err := u.Confirm(fmt.Sprintf("Do you really want to delete the role %s and all of its assignments?", label))
		if err != nil {
			return err
		}
	}

	client, err := provider.Client(c.Context)
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

	return u.Template("The role {{ . | bold }} has been deleted.\n", label.String())
}

func rolesListCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	roles, err := client.ListRoles(c.Context)
	if err != nil {
		return err
	}

	if len(roles) > 0 {
		header := []string{"Label"}
		body := make([][]string, len(roles))
		for i, r := range roles {
			body[i] = []string{r.Label.String()}
		}

		err = u.Table(header, body)
		if err != nil {
			return err
		}
	}

	return u.Template("\nFound {{ . | toString | faded }} role{{ . | pluralize \"s\"}}\n", len(roles))
}

func rolesMembersCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	label := Arguments(c.Context, RoleLabelArg).(primitives.Label)
	client, err := provider.Client(c.Context)
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

	return u.Table(header, body)
}
