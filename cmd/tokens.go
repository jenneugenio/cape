package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func init() {
	tokensCreateCmd := &Command{
		Usage: "Creates a token for the specified identity",
		Examples: []*Example{
			{
				Example:     "cape tokens create",
				Description: "Creates a token for the current user",
			},
			{
				Example:     "cape tokens create service:data-connector@cape.com",
				Description: "Creates a token for the service labelled my-service",
			},
		},
		Arguments: []*Argument{TokenIdentityArg},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(createTokenCmd),
		},
	}

	tokensListCmd := &Command{
		Usage: "Lists the tokens ids for a specified identity",
		Examples: []*Example{
			{
				Example:     "cape tokens list",
				Description: "Lists the token ids for the current user",
			},
			{
				Example:     "cape tokens create service:data-connector@cape.com",
				Description: "Lists the token ids for the service with email service:data-connector@cape.com",
			},
		},
		Arguments: []*Argument{TokenIdentityArg},
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(listTokenCmd),
		},
	}

	tokensRemoveCmd := &Command{
		Usage: "Removes the provided token from Cape",
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
		Usage: "Commands for managing API tokens",
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

func getIdentity(ctx context.Context, client *coordinator.Client) (primitives.Identity, error) {
	var identity primitives.Identity
	identifier, ok := Arguments(ctx, TokenIdentityArg).(primitives.Email)
	if ok {
		identities, err := client.GetIdentities(ctx, []primitives.Email{identifier})
		if err != nil {
			return nil, err
		}

		if len(identities) == 0 {
			return nil, errors.New(NoIdentityCause, "Identity with email %s not found", identifier.String())
		}

		identity = identities[0]
	} else {
		i, err := client.Me(ctx)
		if err != nil {
			return nil, err
		}

		identity = i
	}

	return identity, nil
}

func createTokenCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	identity, err := getIdentity(c.Context, client)
	if err != nil {
		return err
	}

	apiToken, _, err := client.CreateToken(c.Context, identity)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	err = u.Template("A token for {{ . | bold }} has been created!\n\n", identity.GetEmail().String())
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

	identity, err := getIdentity(c.Context, client)
	if err != nil {
		return err
	}

	tokenIDs, err := client.ListTokens(c.Context, identity)
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
