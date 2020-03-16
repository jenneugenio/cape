package cmd

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/logging"
	errors "github.com/dropoutlabs/cape/partyerrors"
)

func init() {
	startCmd := &cli.Command{
		Name:        "start",
		Description: "Starts an instance of the Cape Controlller",
		Action:      startControllerCmd,
		Flags: []cli.Flag{
			dbURLFlag(),
			dbPasswordFlag(),
			instanceIDFlag(),
			loggingTypeFlag(),
			loggingLevelFlag(),
		},
	}

	controllerCmd := &cli.Command{
		Name:        "controller",
		Description: "Commands for launching and managing an instance of Cape's Controller",
		Subcommands: []*cli.Command{startCmd},
	}

	commands = append(commands, controllerCmd)
}

func catchShutdown(ctx context.Context, quit chan os.Signal, c *controller.Controller) error {
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := c.Stop(ctx)
	if err != nil {
		return err
	}

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

	ctx := context.Background()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go catchShutdown(ctx, quit, ctrl) //nolint: errcheck
	return ctrl.Start(ctx)
}

// getDBURL looks at the environment and generates the database address if
// needed.
func getDBURL(c *cli.Context) (*url.URL, error) {
	// We support passing the password in separately or as a part of the DB
	// URL. If the password is contained in the CAPE_DB_URL then it should be
	// passed entirely as a secret inside a kubernetes orchestration system.
	dbURL := c.String("db-url")
	password := c.String("db-password")

	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, errors.Wrap(InvalidURLCause, err)
	}

	// If the password is passed in via environment variables
	// instead of part of the connection string
	if password != "" {
		query := u.Query()
		query.Add("password", password)
		u.RawQuery = query.Encode()
	}

	return u, nil
}
