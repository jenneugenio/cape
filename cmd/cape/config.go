package main

import (
	"fmt"
	"github.com/capeprivacy/cape/models"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/cape/config"
)

func init() {
	viewConfigCmd := &Command{
		Usage:       "List all local configuration settings.",
		Description: "Use this command to list all configuration settings, their defaults, and current values.",
		Command: &cli.Command{
			Name:   "view",
			Action: viewConfig,
		},
	}

	addClusterCmd := &Command{
		Usage:       "Add configuration for a cape cluster.",
		Description: "Use this command to add configuration enabling a user to execute commands against a cluster of Cape.",
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
		Usage:       "Remove configuration for a cape cluster.",
		Description: "Use this command to remove local connection information for a cape cluster.",
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
		Usage:       "Commands for adding, removing, and selecting cape clusters.",
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
		Usage: "Commands for setting and viewing local command line configuration.",
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
	cfg := Config(c.Context)

	label := Arguments(c.Context, ClusterLabelArg).(models.Label)
	clusterURL := Arguments(c.Context, ClusterURLArg).(*models.URL)
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

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)
	err = u.Template("The {{ . | bold }} cluster has been added to your configuration.\n", label.String())
	if err != nil {
		return err
	}

	if use {
		err = u.Template("Your current cluster has been set to {{ . | bold }}.\n", cluster.Label.String())
		if err != nil {
			return err
		}
	}

	return nil
}

func removeCluster(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	cfg := Config(c.Context)
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	// XXX: cluster can be nil here and that's ok - it's possible no cluster
	// was set yet the cluster were trying to delete may exist!
	//
	// We just need to be careful with how we use the `cluster` variable in the
	// rest of the command
	cluster, err := cfg.Cluster()
	if err != nil && err != config.ErrNoCluster {
		return err
	}

	label := Arguments(c.Context, ClusterLabelArg).(models.Label)
	if !skipConfirm {
		// TODO -- template doesn't work with confirm right now
		err := u.Confirm(fmt.Sprintf("Do you want to delete the '%s' cluster from configuration?", label))
		if err != nil {
			return err
		}
	}

	err = cfg.RemoveCluster(label)
	if err != nil {
		return err
	}

	if cluster != nil && cluster.Label == label {
		err = cfg.Use("")
		if err != nil {
			return err
		}
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	err = u.Template("The cluster {{ . | bold }} has been removed from your configuration.\n", label.String())
	if err != nil {
		return err
	}

	if cluster != nil && cluster.Label == label {
		err = u.Template("A current cluster is no longer set. You can set one with {{ . | italic | faded }}\n", "cape config clusters use <label>")
		if err != nil {
			return err
		}
	}

	return nil
}

func useCluster(c *cli.Context) error {
	cfg := Config(c.Context)
	label := Arguments(c.Context, ClusterLabelArg).(models.Label)
	err := cfg.Use(label)
	if err != nil {
		return err
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)
	return u.Template("Your current cluster has been set to {{ . | bold }}.\n", label.String())
}
