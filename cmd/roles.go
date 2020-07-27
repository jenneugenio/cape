package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/primitives"
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
				Example: "cape roles me",
				Description: "Tells you what your role in the org is (admin or user)",
			},
			{
				Example: "cape roles me --project my-project",
				Description: "Tells you what your role in the specified project is",
			},
		},
		Command: &cli.Command{
			Name: "me",
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

	users, err := client.GetUsers(c.Context, emails)
	if err != nil {
		return err
	}

	userIDs := make([]string, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	_, err = client.CreateRole(c.Context, label, userIDs)
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

	users, err := client.GetMembersRole(c.Context, role.ID)
	if err != nil {
		return err
	}

	header := []string{"Email"}
	body := make([][]string, len(users))
	for i, user := range users {
		body[i] = []string{user.Email.String()}
	}

	return u.Table(header, body)
}
