package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	viewConfigCmd := &Command{
		Usage:       "List all local configuration settings",
		Description: "Use this command to list all configuration settings, their defaults, and current values.",
		Command: &cli.Command{
			Name:   "view",
			Action: viewConfig,
		},
	}

	addClusterCmd := &Command{
		Usage:       "Add configuration for a cape cluster",
		Description: "Use this command to add configuration enabling a user to execute commands against a cluster of Cape",
		Arguments:   []*Argument{ClusterLabelArg, ClusterURLArg},
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
		Arguments:   []*Argument{ClusterLabelArg},
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
		Arguments:   []*Argument{ClusterLabelArg},
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
	use := c.Bool("use")
	args := Arguments(c.Context)
	cfg := Config(c.Context)

	label := args["label"].(primitives.Label)
	clusterURL := args["url"].(*url.URL)

	cluster, err := cfg.AddCluster(label, clusterURL, "")
	if err != nil {
		return err
	}

	if use {
		err = cfg.Use(cluster.Label)
		if err != nil {
			return err
		}
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	fmt.Printf("The '%s' cluster has been added to your configuration.\n", label)
	if use {
		fmt.Printf("\nYour current cluster has been set to '%s'.\n", cluster.Label)
	}

	return nil
}

func removeCluster(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	args := Arguments(c.Context)
	cfg := Config(c.Context)

	if skipConfirm {
		fmt.Printf("This is where i should prompt\n")
	}

	cluster, err := cfg.Cluster()
	if err != nil {
		return err
	}

	label := args["label"].(primitives.Label)
	err = cfg.RemoveCluster(label)
	if err != nil {
		return err
	}

	if cluster.Label == label {
		err = cfg.Use(primitives.Label(""))
		if err != nil {
			return err
		}
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	fmt.Printf("The cluster '%s' has been removed from your configation.\n", label)
	if cluster.Label == label {
		fmt.Printf("\nA current cluster is no longer set, please set one using 'cape config clusters use <label>'\n")
	}

	return nil
}

func useCluster(c *cli.Context) error {
	args := Arguments(c.Context)
	cfg := Config(c.Context)

	label := args["label"].(primitives.Label)
	err := cfg.Use(label)
	if err != nil {
		return err
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	fmt.Printf("Your current cluster has been set to '%s'.\n", label)
	return nil
}
