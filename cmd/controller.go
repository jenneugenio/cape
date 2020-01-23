package cmd

import (
	"github.com/dropoutlabs/privacyai/controller"
	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Control access to your data in PrivacyAI",
}

var startControllerCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the PrivacyAI Controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := controller.New()
		c.Start()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)
	controllerCmd.AddCommand(startControllerCmd)
}
