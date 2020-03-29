package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/logging"
	errors "github.com/dropoutlabs/cape/partyerrors"
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
		Arguments: []*Argument{LabelArg("controller"), ClusterURLArg},
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

	err := c.Stop(ctx)
	if err != nil {
		return err
	}

	return nil
}

// getDBURL looks at the environment and generates the database address if
// needed.
func getDBURL(c *cli.Context) (*url.URL, error) {
	// We support passing the password in separately or as a part of the DB
	// URL. If the password is contained in the CAPE_DB_URL then it should be
	// passed entirely as a secret inside a kubernetes orchestration system.
	dbURL := c.String("db-url")
	envVars := EnvVariables(c.Context)
	password, _ := envVars["CAPE_DB_PASSWORD"].(string)

	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, errors.Wrap(InvalidURLCause, err)
	}

	// If the password is passed in via environment variables
	// instead of part of the connection string
	if password != "" {
		u.User = url.UserPassword(u.User.Username(), password)
	}

	return u, nil
}

func setupControllerCmd(c *cli.Context) error {
	args := Arguments(c.Context)
	cfg := Config(c.Context)

	label := args["label"].(primitives.Label)
	clusterURL := args["url"].(*primitives.URL)

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

	fmt.Println("Setup Complete! Welcome to Cape!")
	fmt.Printf("Your current cluster has been set to '%s'.\n", label)

	return nil
}

func startControllerCmd(c *cli.Context) error {
	instanceID, err := getInstanceID(c, "controller")
	if err != nil {
		return err
	}

	dbURL, err := getDBURL(c)
	if err != nil {
		return err
	}

	logger, err := logging.Logger(c.String("logger"), c.String("log-level"), instanceID)
	if err != nil {
		return err
	}

	ctrl, err := controller.New(dbURL, logger, instanceID)
	if err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer signal.Stop(quit)

	go catchShutdown(c.Context, quit, ctrl, logger) //nolint: errcheck
	return ctrl.Start(c.Context)
}
