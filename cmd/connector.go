package main

import (
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/connector"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/logging"
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
				loggingTypeFlag(),
				loggingLevelFlag(),
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
	port := c.Int("port")
	token := EnvVariables(c.Context, capeTokenVar).(*auth.APIToken)

	instanceID, err := getInstanceID(c, "connector")
	if err != nil {
		return err
	}

	cfg := &connector.Config{
		InstanceID: instanceID,
		Port:       port,
		Token:      token,
	}

	// TODO: Consider having the "logger" be configured by the server?
	logger, err := logging.Logger(c.String("logger"), c.String("log-level"), instanceID.String())
	if err != nil {
		return err
	}

	conn, err := connector.New(cfg, logger)
	if err != nil {
		return err
	}

	server, err := framework.NewServer(cfg, conn, logger)
	if err != nil {
		return err
	}

	return server.Start(c.Context)
}
