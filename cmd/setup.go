package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/config"
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func init() {
	setupCmd := &Command{
		Usage: "Setup Cape",
		Description: "Bootstraps the given Cape instance. This will create an admin account for you and give you a token " +
			"that you can use with subsequent commands. You can only run this command ONCE per Cape instance.",
		Examples: []*Example{
			{
				Example:     "cape setup local http://localhost:8080",
				Description: "Initialize an admin account on a local cape instance",
			},
		},
		Arguments: []*Argument{CoordinatorLabelArg, ClusterURLArg},
		Variables: []*EnvVar{capePasswordVar},
		Command: &cli.Command{
			Name:   "setup",
			Action: setupCoordinatorCmd,
			Flags: []cli.Flag{
				nameFlag(),
				emailFlag(),
			},
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

	// Since nothing is set up, ie no clusters exist, we make one here that the provider can fetch
	cluster, err := cfg.AddCluster(label, clusterURL, "")
	if err != nil {
		return err
	}

	s := config.NewSession(cfg, cluster)
	c.Context = context.WithValue(c.Context, SessionContextKey, s)

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

	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	_, err = client.Setup(c.Context, name, models.Email(email.Email), password)
	if err != nil {
		return err
	}

	// Now, log in our admin!
	session, err := client.EmailLogin(c.Context, models.Email(email.Email), password)
	if err != nil {
		return err
	}

	// Now that we're logged in, update the token on our cluster
	cluster.SetToken(session.Token)

	// Set the just setup cluster to the current cluster
	err = cfg.Use(label)
	if err != nil {
		return err
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)

	return u.Template("\nSetup Complete! Welcome to Cape!\n\nYour current cluster has been set to {{ . | bold }}\n", label.String())
}
