package main

import (
	"github.com/manifoldco/go-base64"
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/framework"
	"github.com/capeprivacy/cape/logging"
	"github.com/capeprivacy/cape/primitives"
)

func init() {
	startCmd := &Command{
		Usage:     "Start an instance of the Cape coordinator",
		Variables: []*EnvVar{capeDBPassword, capeDBURL},
		Command: &cli.Command{
			Name:   "start",
			Action: startCoordinatorCmd,
			Flags: []cli.Flag{
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
	url := EnvVariables(c.Context, capeDBURL).(*primitives.DBURL)
	_, set := url.User.Password()
	if set {
		// We do not allow users to set the password in the URL as it is a bad security practice
		return nil, errors.New(PasswordInURLCause, "You cannot set the database password in the URL. Please use CAPE_DB_PASSWORD")
	}
	// If the password is passed in via environment variables
	// instead of part of the connection string.
	//
	// As this env variable is optional we have to check to see if the casting
	// was successful
	password := EnvVariables(c.Context, capeDBPassword).(string)
	url.SetPassword(password)

	return url, nil
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
		return err
	}

	rootKey, err := crypto.GenerateKey()
	if err != nil {
		return err
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
		RootKey:    base64.New(rootKey[:]),
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
