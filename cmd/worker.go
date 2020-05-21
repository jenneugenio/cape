package main

import (
	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/worker"
	"github.com/capeprivacy/cape/primitives"
	"github.com/urfave/cli/v2"
)

func init() {
	startCmd := &Command{
		Usage:     "Start an instance of the cape worker",
		Variables: []*EnvVar{capeTokenVar, capeDBURL},
		Command: &cli.Command{
			Name:   "start",
			Action: startWorkerCmd,
			Flags: []cli.Flag{
				instanceIDFlag(),
				loggingTypeFlag(),
				loggingLevelFlag(),
			},
		},
	}

	workerCmd := &Command{
		Usage: "Commands for managing the Cape worker",
		Command: &cli.Command{
			Name:        "worker",
			Subcommands: []*cli.Command{startCmd.Package()},
		},
	}

	commands = append(commands, workerCmd.Package())
}

func startWorkerCmd(c *cli.Context) error {
	token := EnvVariables(c.Context, capeTokenVar).(*auth.APIToken)
	dbURL := EnvVariables(c.Context, capeDBURL).(*primitives.DBURL)

	config, err := worker.NewConfig(token, dbURL)
	if err != nil {
		return err
	}

	w, err := worker.NewWorker(config)
	if err != nil {
		return err
	}

	return w.Start()
}
