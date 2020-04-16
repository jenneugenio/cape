package main

import (
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/logging"
	"github.com/capeprivacy/cape/primitives"
)

func init() {
	startCmd := &Command{
		Usage:     "Start an instance of the Cape coordinator",
		Variables: []*EnvVar{capeDBPassword},
		Command: &cli.Command{
			Name:   "start",
			Action: startCoordinatorCmd,
			Flags: []cli.Flag{
				dbURLFlag(),
				instanceIDFlag(),
				loggingTypeFlag(),
				loggingLevelFlag(),
				portFlag("coordinator", 8080),
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

// getDBURL looks at the environment and generates the database address if
// needed.
func getDBURL(c *cli.Context) (*primitives.DBURL, error) {
	// We support passing the password in separately or as a part of the DB
	// URL. If the password is contained in the CAPE_DB_URL then it should be
	// passed entirely as a secret inside a kubernetes orchestration system.
	u, err := primitives.NewDBURL(c.String("db-url"))
	if err != nil {
		return nil, err
	}

	// If the password is passed in via environment variables
	// instead of part of the connection string.
	//
	// As this env variable is optional we have to check to see if the casting
	// was successful
	password, ok := EnvVariables(c.Context, capeDBPassword).(string)
	if ok && password != "" {
		u.SetPassword(password)
	}

	return u, nil
}

func startCoordinatorCmd(c *cli.Context) error {
	port := c.Int("port")
	instanceID, err := getInstanceID(c, "coordinator")
	if err != nil {
		return err
	}

	dbURL, err := getDBURL(c)
	if err != nil {
		return err
	}

	// TODO: Consider having the "logger" be configured by the server?
	logger, err := logging.Logger(c.String("logger"), c.String("log-level"), instanceID.String())
	if err != nil {
		return err
	}

	// TODO: Finish integrating loading config from a file and enabling
	// overwriting of flags. This includes figuring out the local development
	// environment and configuration workflow.
	keypair, err := auth.NewKeypair()
	if err != nil {
		return nil
	}

	cfg := &coordinator.Config{
		DB: &coordinator.DBConfig{
			Addr: dbURL,
		},
		Auth: &coordinator.AuthConfig{
			KeypairPackage: keypair.Package(),
		},
		InstanceID: instanceID,
		Port:       port,
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
