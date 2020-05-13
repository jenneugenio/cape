package main

import (
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/coordinator/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

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

func updateCmd(c *cli.Context) error {
	dbURL, err := getDBURL(c)
	if err != nil {
		return err
	}

	migrator, err := database.NewMigrator(dbURL.URL, "coordinator/migrations")
	if err != nil {
		return err
	}

	return migrator.Up(c.Context)
}

func init() {
	updateCmd := &Command{
		Usage:     "Update a running Cape coordinator to a new version",
		Variables: []*EnvVar{capeDBPassword, capeDBURL},
		Command: &cli.Command{
			Name:   "update",
			Action: updateCmd,
		},
	}

	commands = append(commands, updateCmd.Package())
}
