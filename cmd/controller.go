package cmd

import (
	"net/url"

	"github.com/dropoutlabs/privacyai/controller"
	errors "github.com/dropoutlabs/privacyai/partyerrors"
	"github.com/urfave/cli/v2"
)

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

func startControllerCmd(c *cli.Context) error {
	serviceID := getServiceID(c)
	dbURL, err := getDBURL(c)
	if err != nil {
		return err
	}

	ctrl, err := controller.New(dbURL, serviceID)
	if err != nil {
		return err
	}

	ctrl.Start()

	return nil
}

func init() {
	startCmd := &cli.Command{
		Name:        "start",
		Description: "Launch the Cape Controller",
		Action:      startControllerCmd,
		Flags: []cli.Flag{
			dbURLFlag(),
			dbPasswordFlag(),
			serviceIDFlag(),
		},
	}

	controllerCmd := &cli.Command{
		Name:        "controller",
		Description: "Control access to your data in PrivacyAI",
		Subcommands: []*cli.Command{startCmd},
	}

	commands = append(commands, controllerCmd)
}
