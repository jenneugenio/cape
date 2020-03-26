package main

import (
	"crypto/rand"
	"fmt"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
	"github.com/manifoldco/go-base64"
	"github.com/urfave/cli/v2"
)

const (
	secretBytesLength = 8
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
			Action: usersCreateCmd,
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
	ui := UI(c.Context)
	args := Arguments(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	URL, err := cluster.GetURL()
	if err != nil {
		return err
	}

	token, err := cluster.Token()
	if err != nil {
		return err
	}

	email := args["email"].(primitives.Email)

	validateName := func(input string) error {
		_, err := primitives.NewName(input)
		if err != nil {
			return err
		}

		return nil
	}

	nameStr, err := ui.Question("Name", validateName)
	if err != nil {
		return err
	}

	name, err := primitives.NewName(nameStr)
	if err != nil {
		return err
	}

	secretBytes := make([]byte, secretBytesLength)
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

	client := controller.NewClient(URL, token)
	_, err = client.CreateUser(c.Context, user)
	if err != nil {
		return err
	}

	fmt.Println("A user has been created with the following credentials:")
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Emai: %s\n", email)
	fmt.Printf("Passowrd %s\n", password)
	fmt.Println("Remember: Please keep the password safe and share with the user securely.")

	return nil
}
