package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func init() {
	logoutCmd := &Command{
		Usage: "Removes session on the client and the server",
		Command: &cli.Command{
			Name:   "logout",
			Action: logoutCmd,
			Flags: []cli.Flag{
				emailFlag(),
				passwordFlag(),
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

	client, err := cluster.Client()
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

	fmt.Printf("You have been logged out of '%s'.\n", cluster.String())
	return nil
}
