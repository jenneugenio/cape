package main

import (
	"fmt"
	"github.com/capeprivacy/cape/auth"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
	"github.com/urfave/cli/v2"
)

func init() {
	setupCmd := &Command{
		Usage: "Setup Cape",
		Description: "Bootstraps the given Cape instance. This will create an admin account for you and give you a token " +
			"that you can use with subsequent commands. You can only run this command ONCE per Cape instance.",
		Examples: []*Example{
			{
				Example:     "cape coordinator setup local http://localhost:8081",
				Description: "Initialize an admin account on a local cape instance",
			},
		},
		Arguments: []*Argument{CoordinatorLabelArg, ClusterURLArg},
		Command: &cli.Command{
			Name:   "setup",
			Action: setupCoordinatorCmd,
		},
	}

	commands = append(commands, setupCmd.Package())
}

func setupCoordinatorCmd(c *cli.Context) error {
	cfg := Config(c.Context)

	label := Arguments(c.Context, CoordinatorLabelArg).(primitives.Label)
	clusterURL := Arguments(c.Context, ClusterURLArg).(*primitives.URL)

	if cfg.HasCluster(label) {
		removeCmd := fmt.Sprintf("cape config clusters remove %s", label)
		return errors.New(ClusterExistsCause, "A cluster named '%s' has already been configured! You can use `%s` to remove it.", label, removeCmd)
	}

	name, err := getName(c, "")
	if err != nil {
		return err
	}

	email, err := getEmail(c, "")
	if err != nil {
		return err
	}

	password, err := getConfirmedPassword(c)
	if err != nil {
		return err
	}

	creds, err := auth.NewCredentials(password.Bytes(), nil)
	if err != nil {
		return err
	}

	user, err := primitives.NewUser(name, email, creds.Package())
	if err != nil {
		return err
	}

	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	_, err = client.Setup(c.Context, user)
	if err != nil {
		return err
	}

	// Now, log in our admin!
	session, err := client.Login(c.Context, email, []byte(password))
	if err != nil {
		return err
	}

	_, err = cfg.AddCluster(label, clusterURL, session.Token.String())
	if err != nil {
		return err
	}

	err = cfg.Use(label)
	if err != nil {
		return err
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	fmt.Printf("\nSetup Complete! Welcome to Cape!\n\n")
	fmt.Printf("Your current cluster has been set to '%s'.\n", label)

	return nil
}
