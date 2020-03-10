package cmd

import (
	"github.com/dropoutlabs/cape/connector"
	"github.com/urfave/cli/v2"
)

func startConnectorCmd(c *cli.Context) error {
	serviceID := getServiceID(c)
	conn := connector.New(serviceID)
	conn.Start()

	return nil
}

func init() {
	startCmd := &cli.Command{
		Name:        "start",
		Description: "Launch the Cape Data Connector",
		Action:      startConnectorCmd,
		Flags:       []cli.Flag{serviceIDFlag()},
	}

	connectorCmd := &cli.Command{
		Name:        "connector",
		Description: "Connect your data sources for use within Cape",
		Subcommands: []*cli.Command{startCmd},
	}

	commands = append(commands, connectorCmd)
}
