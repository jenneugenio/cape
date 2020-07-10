package main

import (
	"github.com/capeprivacy/cape/models"
	"github.com/urfave/cli/v2"
)

func init() {
	loginCmd := &Command{

		Usage:     "Creates a session on the coordinator",
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

	email, err := getEmail(c, c.String("email"))
	if err != nil {
		return err
	}
	password, err := getPassword(c)
	if err != nil {
		return err
	}

	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	session, err := client.EmailLogin(c.Context, models.Email(email.Email), password)
	if err != nil {
		return err
	}

	cluster.SetToken(session.Token)
	err = cfg.Write()
	if err != nil {
		return err
	}

	args := struct {
		Email      string
		ClusterURL string
	}{
		email.String(),
		cluster.URL.String(),
	}

	u := provider.UI(c.Context)
	return u.Template("You are now authenticated to {{ .ClusterURL | bold }} as {{ .Email | bold }}\n", args)
}
