package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/framework"
	"github.com/dropoutlabs/cape/logging"
	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	setupCmd := &Command{
		Usage: "Setup an instance of the Cape controller",
		Description: "Bootstraps the given Cape instance. This will create an admin account for you and give you a token " +
			"that you can use with subsequent commands. You can only run this command ONCE per Cape instance.",
		Examples: []*Example{
			{
				Example:     "cape controller setup local http://localhost:8081",
				Description: "Initialize an admin account on a local cape instance",
			},
		},
		Arguments: []*Argument{ControllerLabelArg, ClusterURLArg},
		Command: &cli.Command{
			Name:   "setup",
			Action: setupControllerCmd,
		},
	}

	startCmd := &Command{
		Usage:     "Start an instance of the Cape controller",
		Variables: []*EnvVar{capeDBPassword},
		Command: &cli.Command{
			Name:   "start",
			Action: startControllerCmd,
			Flags: []cli.Flag{
				dbURLFlag(),
				instanceIDFlag(),
				loggingTypeFlag(),
				loggingLevelFlag(),
				portFlag("controller", 8080),
			},
		},
	}

	controllerCmd := &Command{
		Usage: "Commands for starting and managing Cape controllers.",
		Command: &cli.Command{
			Name:        "controller",
			Subcommands: []*cli.Command{setupCmd.Package(), startCmd.Package()},
		},
	}

	commands = append(commands, controllerCmd.Package())
}

func catchShutdown(ctx context.Context, quit chan os.Signal, c *controller.Controller, logger *zerolog.Logger) error {
	s := <-quit

	// Once we've received a signal we will stop listening for additional
	// signals. This will cause the behaviour to fall back to the underlying
	// golang signal handler which will cause in the program immediately
	// exiting.
	//
	// This is not desirable as we will not clean up all of the state but its
	// the best option if we get stuck in an irrecoverable state (e.g. things
	// are not timing out).
	signal.Stop(quit)

	logger.Warn().Msgf("Caught signal '%s': attempting to shutdown, 30 second timeout.", s.String())
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := c.Teardown(ctx)
	if err != nil {
		return err
	}

	return nil
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

func setupControllerCmd(c *cli.Context) error {
	cfg := Config(c.Context)

	label := Arguments(c.Context, ControllerLabelArg).(primitives.Label)
	clusterURL := Arguments(c.Context, ClusterURLArg).(*primitives.URL)
	name, err := getName(c, "")
	if err != nil {
		return err
	}

	email, err := getEmail(c, "")
	if err != nil {
		return err
	}

	password, err := getConfirmedPassword(c)
	if err != nil {
		return err
	}

	creds, err := auth.NewCredentials(password.Bytes(), nil)
	if err != nil {
		return err
	}

	user, err := primitives.NewUser(name, email, creds.Package())
	if err != nil {
		return err
	}

	client := controller.NewClient(clusterURL, nil)
	_, err = client.Setup(c.Context, user)
	if err != nil {
		return err
	}

	// Now, log in our admin!
	session, err := client.Login(c.Context, email, []byte(password))
	if err != nil {
		return err
	}

	_, err = cfg.AddCluster(label, clusterURL, session.Token.String())
	if err != nil {
		return err
	}

	err = cfg.Use(label)
	if err != nil {
		return err
	}

	err = cfg.Write()
	if err != nil {
		return err
	}

	fmt.Printf("\nSetup Complete! Welcome to Cape!\n\n")
	fmt.Printf("Your current cluster has been set to '%s'.\n", label)

	return nil
}

func startControllerCmd(c *cli.Context) error {
	port := c.Int("port")
	instanceID, err := getInstanceID(c, "controller")
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

	cfg := &controller.Config{
		DB: &controller.DBConfig{
			Addr: dbURL,
		},
		Auth: &controller.AuthConfig{
			KeypairPackage: keypair.Package(),
		},
		InstanceID: instanceID,
		Port:       port,
	}

	ctrl, err := controller.New(cfg, logger)
	if err != nil {
		return err
	}

	server, err := framework.NewServer(cfg, ctrl, logger)
	if err != nil {
		return err
	}

	// TODO: Move all of the signal bits into framework.Server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer signal.Stop(quit)

	go catchShutdown(c.Context, quit, ctrl, logger) //nolint: errcheck
	return server.Start(c.Context)
}
