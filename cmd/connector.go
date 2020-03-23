package main

import (
	"github.com/dropoutlabs/cape/connector"
	"github.com/urfave/cli/v2"
)

func startConnectorCmd(c *cli.Context) error {
	instanceID, err := getInstanceID(c, "connector")
	if err != nil {
		return err
	}

	conn := connector.New(instanceID)
	conn.Start()

	return nil
}

func init() {
	startCmd := &Command{
		Usage:       "Start an instance of the Cape data connector",
		Description: "Use this command to start an instance of a Cape data connector.",
		Command: &cli.Command{
			Name:   "start",
			Action: startConnectorCmd,
			Flags:  []cli.Flag{instanceIDFlag()},
		},
	}

	connectorCmd := &Command{
		Usage:       "Commands for starting and managing data connectors",
		Description: "Commands for managing Cape data connectors which make data sources available for use within Cape.",
		Command: &cli.Command{
			Name:        "connector",
			Subcommands: []*cli.Command{startCmd.Package()},
		},
	}

	commands = append(commands, connectorCmd.Package())
}
