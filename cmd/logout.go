package main

import (
	"github.com/urfave/cli/v2"
)

func init() {
	logoutCmd := &Command{
		Usage: "Removes session on the client and the server.",
		Command: &cli.Command{
			Name:   "logout",
			Action: handleSessionOverrides(logoutCmd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	commands = append(commands, logoutCmd.Package())
}

func logoutCmd(c *cli.Context) error {
	cfgSession := Session(c.Context)
	cfg := Config(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	err = client.Logout(c.Context, nil)
	if err != nil {
		return err
	}

	cluster.SetToken(nil)
	err = cfg.Write()
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("You have been logged out of {{ . | bold }}\n", cluster.String())
}
