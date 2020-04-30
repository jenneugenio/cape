package main

import (
	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/ui"
)

func init() {
	tokensCreateCmd := &Command{
		Usage: "Creates a token for the specified entity",
		Examples: []*Example{
			{
				Example:     "cape tokens create",
				Description: "Creates a token for the current user",
			},
			{
				Example:     "cape tokens create my-service",
				Description: "Creates a token for the service labelled my-service",
			},
		},
		Arguments: []*Argument{TokenIdentityArg},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(createTokenCmd),
		},
	}

	tokensCmd := &Command{
		Usage: "Commands for managing API tokens",
		Command: &cli.Command{
			Name: "tokens",
			Subcommands: []*cli.Command{
				tokensCreateCmd.Package(),
			},
		},
	}

	commands = append(commands, tokensCmd.Package())
}

func createTokenCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	token, err := client.CreateToken(c.Context, nil)
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)
	err = u.Template("A token has been created!\n\n", nil)
	if err != nil {
		return err
	}

	tokenStr, err := token.Marshal()
	if err != nil {
		return err
	}

	err = u.Details(ui.Details{"Token": tokenStr})
	if err != nil {
		return err
	}

	return u.Notify(ui.Remember, "Please keep the token safe and share it only over secure channels.")
}
