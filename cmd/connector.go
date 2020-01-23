package cmd

import (
	"github.com/dropoutlabs/privacyai/connector"
	"github.com/spf13/cobra"
)

var connectorCmd = &cobra.Command{
	Use:   "connector",
	Short: "Connect your data sources for use within PrivacyAI",
}

var startConnectorCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the PrivacyAI Data Connector",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := connector.New()
		c.Start()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectorCmd)
	connectorCmd.AddCommand(startConnectorCmd)
}
