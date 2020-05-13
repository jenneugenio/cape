package main

import (
	"io/ioutil"

	"github.com/urfave/cli/v2"
	"sigs.k8s.io/yaml"

	"github.com/kelseyhightower/envconfig"

	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/logging"
)

func init() {
	startCmd := &Command{
		Usage: "Start an instance of the Cape coordinator",
		Command: &cli.Command{
			Name:   "start",
			Action: startCoordinatorCmd,
			Flags: []cli.Flag{
				loggingTypeFlag(),
				loggingLevelFlag(),
				configFileFlag(),
			},
		},
	}

	coordinatorCmd := &Command{
		Usage: "Commands for starting and managing Cape coordinators.",
		Command: &cli.Command{
			Name:        "coordinator",
			Subcommands: []*cli.Command{startCmd.Package()},
		},
	}

	commands = append(commands, coordinatorCmd.Package())
}

func getConfig(c *cli.Context) (*coordinator.Config, error) {
	configs := c.StringSlice("file")

	config := &coordinator.Config{}
	for _, configFile := range configs {
		by, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(by, config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func startCoordinatorCmd(c *cli.Context) error {
	cfg, err := getConfig(c)
	if err != nil {
		return err
	}

	err = envconfig.Process("cape", cfg)
	if err != nil {
		return err
	}

	// TODO: Consider having the "logger" be configured by the server?
	logger, err := logging.Logger(c.String("logger"), c.String("log-level"), cfg.InstanceID.String())
	if err != nil {
		return err
	}

	ctrl, err := coordinator.New(cfg, logger)
	if err != nil {
		return err
	}

	server, err := framework.NewServer(cfg, ctrl, logger)
	if err != nil {
		return err
	}

	watcher, err := setupSignalWatcher(server, logger)
	if err != nil {
		return err
	}

	err = watcher.Start()
	if err != nil {
		return err
	}
	defer watcher.Stop()

	return server.Start(c.Context)
}
