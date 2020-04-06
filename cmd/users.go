package main

import (
	"crypto/rand"
	"fmt"

	"github.com/manifoldco/go-base64"
	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/cmd/ui"
	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	createCmd := &Command{
		Usage:     "Create a new user",
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
		Usage: "Commands for querying information about users and modifying them",
		Command: &cli.Command{
			Name:        "users",
			Subcommands: []*cli.Command{createCmd.Package()},
		},
	}

	commands = append(commands, usersCmd.Package())
}

func usersCreateCmd(c *cli.Context) error {
	args := Arguments(c.Context)
	u := UI(c.Context)
	cfgSession := Session(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	email := args["email"].(primitives.Email)
	name, err := getName(c, "Please enter the persons name")
	if err != nil {
		return err
	}

	secretBytes := make([]byte, auth.GeneratedSecretByteLength)
	_, err = rand.Read(secretBytes)
	if err != nil {
		return err
	}

	password := base64.New(secretBytes).String()
	creds, err := auth.NewCredentials([]byte(password), nil)
	if err != nil {
		return err
	}

	user, err := primitives.NewUser(name, email, creds.Package())
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	_, err = client.CreateUser(c.Context, user)
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

	return u.Notify(ui.Remember, "Please keep the password safe and share it only over secure channels.\n")
}
