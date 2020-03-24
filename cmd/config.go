package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func init() {
	viewConfigCmd := &Command{
		Usage:       "List all local configuration settings",
		Description: "Use this command to list all configuration settings, their defaults, and current values.",
		Command: &cli.Command{
			Name:   "view",
			Action: retrieveConfig(viewConfig),
		},
	}

	labelArg := NewArgument("label", "A label for the cluster", true)
	urlArg := NewArgument("url", "The url of the cape cluster", true)
	addClusterCmd := &Command{
		Usage:       "Add configuration for a cape cluster",
		Description: "Use this command to add configuration enabling a user to execute commands against a cluster of Cape",
		Arguments:   []*Argument{labelArg, urlArg},
		Examples: []*Example{
			{
				Example:     "cape config clusters add production https://my.production.com",
				Description: "Add a cluster labeled 'production'.",
			},
			{
				Example:     "cape config clusters add --use production https://my.production.com",
				Description: "Add a cluster named 'production' and switch to using it.",
			},
		},
		Command: &cli.Command{
			Name:   "add",
			Flags:  []cli.Flag{useClusterFlag()},
			Action: addCluster,
		},
	}

	removeClusterCmd := &Command{
		Usage:       "Remove configuration for a cape cluster",
		Description: "Use this command to remove local connection information for a cape cluster",
		Arguments:   []*Argument{labelArg},
		Examples: []*Example{
			{
				Example:     "cape config clusters remove production",
				Description: "Remove all configuration related to a cluster named 'production'",
			},
			{
				Example:     "cape config clusters remove production -y",
				Description: "Remove all configuration related to a cluster named 'production' without prompting for confirmation.",
			},
		},
		Command: &cli.Command{
			Name:   "remove",
			Action: removeCluster,
			Flags:  []cli.Flag{yesFlag()},
		},
	}

	useDescription := "Use this comman to set a configured cape cluster as the current cluster. " +
		"All subsequent commands will be executed against this cluster."
	useClusterCmd := &Command{
		Usage:       "Set a cape cluster as your current cluster",
		Description: useDescription,
		Arguments:   []*Argument{labelArg},
		Examples: []*Example{
			{
				Example:     "cape config clusters use production",
				Description: "Run all commands against the 'production' cluster",
			},
		},
		Command: &cli.Command{
			Name:   "use",
			Action: useCluster,
		},
	}

	clustersCmd := &Command{
		Usage:       "Commands for adding, removing, and selecting cape clusters",
		Description: "Use these commands for adding, removing, and selecting your current cape cluster.",
		Command: &cli.Command{
			Name: "clusters",
			Subcommands: []*cli.Command{
				addClusterCmd.Package(),
				removeClusterCmd.Package(),
				useClusterCmd.Package(),
			},
		},
	}

	configCmd := &Command{
		Usage: "Commands for setting and viewing local command line configuration",
		Command: &cli.Command{
			Name: "config",
			Subcommands: []*cli.Command{
				viewConfigCmd.Package(),
				clustersCmd.Package(),
			},
		},
	}

	commands = append(commands, configCmd.Package())
}

func viewConfig(c *cli.Context) error {
	cfg := Config(c.Context)
	return cfg.Print(os.Stdout)
}

func addCluster(c *cli.Context) error {
	return nil
}

func removeCluster(c *cli.Context) error {
	return nil
}

func useCluster(c *cli.Context) error {
	return nil
}
