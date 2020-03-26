package main

import (
	"fmt"

	"github.com/dropoutlabs/cape/auth"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/controller"
)

func init() {
	loginCmd := &Command{
		Usage: "Creates a session on the controller",
		Command: &cli.Command{
			Name:   "login",
			Action: loginCmd,
			Flags: []cli.Flag{
				emailFlag(),
				passwordFlag(),
				clusterFlag(),
			},
		},
	}

	commands = append(commands, loginCmd.Package())
}

func loginCmd(c *cli.Context) error {
	ui := UI(c.Context)
	cfg := Config(c.Context)
	cfgSession := Session(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	emailStr := c.String("email")
	password := c.String("password")

	if emailStr == "" {
		tmpE, err := ui.Question("Please enter your email address", nil)
		if err != nil {
			return err
		}
		emailStr = tmpE
	}

	email, err := primitives.NewEmail(emailStr)
	if err != nil {
		return err
	}

	validatePassword := func(input string) error {
		if len(input) < auth.SecretLength {
			return errors.New(InvalidLengthCause, "Password is too short")
		}

		return nil
	}

	if password == "" {
		password, err = ui.Secret("Please enter a password", validatePassword)
		if err != nil {
			return err
		}
	} else {
		err = validatePassword(password)
		if err != nil {
			return err
		}
	}

	URL, err := cluster.GetURL()
	if err != nil {
		return err
	}

	client := controller.NewClient(URL, nil)
	session, err := client.Login(c.Context, email, []byte(password))
	if err != nil {
		return err
	}

	cluster.SetToken(session.Token)

	err = cfg.Write()
	if err != nil {
		return err
	}

	fmt.Printf("You are now authenticated as %s to '%s'.\n", email, cluster.String())

	return nil
}
