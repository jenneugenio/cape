package cmd

import (
	"net/url"
	"os"

	"github.com/dropoutlabs/privacyai/controller"
	errors "github.com/dropoutlabs/privacyai/partyerrors"
	"github.com/spf13/cobra"
)

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Control access to your data in PrivacyAI",
}

// getDBURL looks at the environment and generates the database address if
// needed.
func getDBURL() (*url.URL, error) {
	// We support passing the password in separately or as a part of the DB
	// URL. If the password is contained in the DB_URL then it should be passed
	// entirely as a secret inside a kubernetes orchestration system.
	dbURL := os.Getenv("DB_URL")
	password := os.Getenv("DB_PASSWORD")

	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, errors.Wrap(InvalidURLCause, err)
	}

	// If the password is passed in via and environment variables
	// instead of part of the connection string
	if password != "" {
		query := u.Query()
		query.Add("password", password)
		u.RawQuery = query.Encode()
	}

	return u, nil
}

var startControllerCmd = &cobra.Command{
	Use:   "start",
	Short: "Launch the PrivacyAI Controller",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbURL, err := getDBURL()
		if err != nil {
			return err
		}

		c, err := controller.New(dbURL)
		if err != nil {
			return err
		}

		c.Start()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(controllerCmd)
	controllerCmd.AddCommand(startControllerCmd)
}
