package main

import (
	"github.com/dropoutlabs/cape/framework"
	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/connector"
)

func init() {
	startCmd := &Command{
		Usage:       "Start an instance of the Cape data connector",
		Description: "Use this command to start an instance of a Cape data connector.",
		Variables:   []*EnvVar{capeTokenVar},
		Command: &cli.Command{
			Name:   "start",
			Action: startConnectorCmd,
			Flags: []cli.Flag{
				instanceIDFlag(),
				portFlag("connector", 8081),
			},
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

func startConnectorCmd(c *cli.Context) error {
	envVars := EnvVariables(c.Context)
	port := c.Int("port")
	token := envVars["CAPE_TOKEN"].(*auth.APIToken)

	instanceID, err := getInstanceID(c, "connector")
	if err != nil {
		return err
	}

	cfg := &connector.Config{
		InstanceID: instanceID,
		Port:       port,
		Token:      token,
	}

	conn, err := connector.New(cfg)
	if err != nil {
		return err
	}

	server, err := framework.NewServer(cfg, conn)
	if err != nil {
		return err
	}

	return server.Start(c.Context)
}
