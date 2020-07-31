package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/cape/ui"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func init() {
	createCmd := &Command{
		Usage:     "Create a new user.",
		Arguments: []*Argument{UserEmailArg},
		Examples: []*Example{
			{
				Example: "cape users create email@email.com",
				Description: "Creates a new user with the email 'email@email.com'. Will prompt " +
					"for name and generate a password.",
			},
			{
				Example:     "CAPE_CLUSTER=prod cape users create email@email.com",
				Description: "Creates a user on the cape cluster called prod.",
			},
		},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(usersCreateCmd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	usersCmd := &Command{
		Usage: "Commands for querying information about users and modifying them.",
		Command: &cli.Command{
			Name:        "users",
			Subcommands: []*cli.Command{createCmd.Package()},
		},
	}

	commands = append(commands, usersCmd.Package())
}

func usersCreateCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	email := Arguments(c.Context, UserEmailArg).(primitives.Email)
	name, err := getName(c, "Please enter the persons name")
	if err != nil {
		return err
	}

	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	_, password, err := client.CreateUser(c.Context, name, models.Email(email.Email))
	if err != nil {
		return err
	}

	fmt.Printf("A user has been created with the following credentials:\n\n")

	err = u.Details(ui.Details{
		"Name":     name,
		"Email":    email,
		"Password": password,
	})
	if err != nil {
		return err
	}

	return u.Notify(ui.Remember, "Please keep the password safe and share it only over secure channels.")
}
