package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	loginCmd := &Command{

		Usage:     "Creates a session on the controller",
		Variables: []*EnvVar{capePasswordVar},
		Command: &cli.Command{
			Name:   "login",
			Action: handleSessionOverrides(loginCmd),
			Flags: []cli.Flag{
				emailFlag(),
				clusterFlag(),
			},
		},
	}

	commands = append(commands, loginCmd.Package())
}

func loginCmd(c *cli.Context) error {
	cfg := Config(c.Context)
	cfgSession := Session(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	email, err := getEmail(c)
	if err != nil {
		return err
	}
	password, err := getPassword(c)
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	session, err := client.Login(c.Context, email, []byte(string(password)))
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

func getEmail(c *cli.Context) (primitives.Email, error) {
	in := c.String("email")
	if in != "" {
		return primitives.NewEmail(in)
	}

	ui := UI(c.Context)
	out, err := ui.Question("Please enter your email address", func(input string) error {
		_, err := primitives.NewEmail(input)
		return err
	})
	if err != nil {
		return primitives.Email(""), err
	}

	return primitives.NewEmail(out)
}

func getPassword(c *cli.Context) (primitives.Password, error) {
	envVars := EnvVariables(c.Context)
	ui := UI(c.Context)

	pw, ok := envVars["CAPE_PASSWORD"].(primitives.Password)
	if ok {
		return pw, nil
	}

	// XXX: It'd be nice if we didn't need to do this weird type creation
	// manipulation. If we could just reuse the `.Validate()` function that'd
	// be awesome butthat's not how the promptui ValidatorFunc works!
	out, err := ui.Secret("Please enter a password", func(input string) error {
		_, err := primitives.NewPassword(input)
		return err
	})
	if err != nil {
		return pw, err
	}

	return primitives.NewPassword(out)
}
