package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/models"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func init() {
	tokensCreateCmd := &Command{
		Usage: "Creates a token for the specified user.",
		Examples: []*Example{
			{
				Example:     "cape tokens create",
				Description: "Creates a token for the current user.",
			},
			{
				Example:     "cape tokens create user@cape.com",
				Description: "Creates a token for the user with the email user@cape.com.",
			},
		},
		Arguments: []*Argument{TokenUserArg},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(createTokenCmd),
		},
	}

	tokensListCmd := &Command{
		Usage: "Lists the tokens IDs for a specified user.",
		Examples: []*Example{
			{
				Example:     "cape tokens list",
				Description: "Lists the token ids for the current user.",
			},
			{
				Example:     "cape tokens create data-user@cape.com",
				Description: "Lists the token ids for the user with email data-user@cape.com.",
			},
		},
		Arguments: []*Argument{TokenUserArg},
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(listTokenCmd),
		},
	}

	tokensRemoveCmd := &Command{
		Usage: "Removes the provided token from Cape.",
		Examples: []*Example{
			{
				Example:     "cape tokens remove 2011e949qta0quff3n4yx7ny3r",
				Description: "Removes the provided token (by ID) from the database",
			},
		},
		Arguments: []*Argument{TokenIDArg},
		Command: &cli.Command{
			Name:   "remove",
			Action: handleSessionOverrides(removeTokenCmd),
		},
	}

	tokensCmd := &Command{
		Usage: "Commands for managing API tokens.",
		Command: &cli.Command{
			Name: "tokens",
			Subcommands: []*cli.Command{
				tokensCreateCmd.Package(),
				tokensListCmd.Package(),
				tokensRemoveCmd.Package(),
			},
		},
	}

	commands = append(commands, tokensCmd.Package())
}

func getUser(ctx context.Context, client *coordinator.Client) (*models.User, error) {
	var user *models.User
	identifier, ok := Arguments(ctx, TokenUserArg).(primitives.Email)
	if ok {
		users, err := client.GetUsers(ctx, []primitives.Email{identifier})
		if err != nil {
			return nil, err
		}

		if len(users) == 0 {
			return nil, errors.New(NoUserCause, "User with email %s not found", identifier.String())
		}

		user = users[0]
	} else {
		i, err := client.Me(ctx)
		if err != nil {
			return nil, err
		}

		user = i
	}

	return user, nil
}

func createTokenCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	user, err := getUser(c.Context, client)
	if err != nil {
		return err
	}

	apiToken, _, err := client.CreateToken(c.Context, user)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	err = u.Template("A token for {{ . | bold }} has been created!\n\n", user.Email.String())
	if err != nil {
		return err
	}

	tokenStr, err := apiToken.Marshal()
	if err != nil {
		return err
	}

	err = u.Details(ui.Details{"Token": tokenStr})
	if err != nil {
		return err
	}

	return u.Notify(ui.Remember, "Please keep the token safe and share it only over secure channels.")
}

func listTokenCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	user, err := getUser(c.Context, client)
	if err != nil {
		return err
	}

	tokenIDs, err := client.ListTokens(c.Context, user)
	if err != nil {
		return err
	}

	header := []string{"Token ID"}
	body := make([][]string, len(tokenIDs))

	for i, t := range tokenIDs {
		body[i] = []string{t.String()}
	}

	u := provider.UI(c.Context)
	err = u.Table(header, body)
	if err != nil {
		return err
	}

	return u.Template("\nFound {{ . | toString | faded }} token{{ . | pluralize \"s\"}}\n", len(tokenIDs))
}

func removeTokenCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	ID, _ := Arguments(c.Context, TokenIDArg).(database.ID)
	err = client.RemoveToken(c.Context, ID)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	return u.Template("Removed the token with ID {{ . | toString | faded }} from Cape\n", ID.String())
}
