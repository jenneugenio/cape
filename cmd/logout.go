package main

import (
	"fmt"
	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/controller"
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

	URL, err := cluster.GetURL()
	if err != nil {
		return err
	}

	token, err := cluster.Token()
	if err != nil {
		return err
	}

	client := controller.NewClient(URL, token)
	err = client.Logout(c.Context, token)
	if err != nil {
		return err
	}

	cluster.AuthToken = ""

	err = cfg.Write()
	if err != nil {
		return err
	}

	fmt.Printf("You have been logged out of '%s'.\n", cluster.String())

	return nil
}
