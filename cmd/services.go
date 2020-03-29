package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	createCmd := &Command{
		Usage:     "Create a new service",
		Arguments: []*Argument{ServiceIdentifierArg},
		Examples: []*Example{
			{
				Example:     "cape services create service:pipeline@prod.mycompany.com",
				Description: "Creates a new service with the email 'pipeline@prod.mycompany.com'.",
			},
			{
				Example:     "CAPE_CLUSTER=production cape services create service:pipeline@prod.mycompany.com",
				Description: "Creates a service called service:pipeline@prod.mycompany.com on the cape cluster called production.",
			},
			{
				Example:     "cape services create --type data-connector service:pipeline@prod.mycompany.com",
				Description: "Creates a service called service:pipeline@prod.mycompany.com representing a Cape data connector.",
			},
			{
				Example: "cape services create pipeline@prod.mycompany.com",
				Description: "Creates a new service with the email 'service:pipeline@prod.mycompany.com'.\n" +
					"'service:' is prepended for you if its not included in the given email",
			},
		},
		Command: &cli.Command{
			Name:   "create",
			Action: handleSessionOverrides(servicesCreateCmd),
			Flags: []cli.Flag{
				clusterFlag(),
				serviceTypeFlag(),
			},
		},
	}

	removeCmd := &Command{
		Usage:     "Remove command removes a service",
		Arguments: []*Argument{ServiceIdentifierArg},
		Examples: []*Example{
			{
				Example:     "cape services remove servce:pipeline@prod.mycompany.com",
				Description: "Removes a new service with the email 'service:pipeline@prod.mycompany.com'.",
			},
			{
				Example:     "cape services remove --yes service:pipeline@prod.mycompany.com",
				Description: "Removes a service skipping the confirm dialog.",
			},
			{
				Example: "cape services remove pipeline@prod.mycompany.com",
				Description: "Removes a service skipping the confirm dialog. " +
					"'service:' is prepended for you if its not included in the given email",
			},
		},
		Command: &cli.Command{
			Name:   "remove",
			Action: handleSessionOverrides(servicesRemoveCmd),
			Flags: []cli.Flag{
				clusterFlag(),
				yesFlag(),
			},
		},
	}

	listCmd := &Command{
		Usage: "Lists all the services on the cluster",
		Examples: []*Example{
			{
				Example:     "cape services list",
				Description: "Lists all services",
			},
		},
		Command: &cli.Command{
			Name:   "list",
			Action: handleSessionOverrides(servicesListCmd),
			Flags: []cli.Flag{
				clusterFlag(),
			},
		},
	}

	servicesCmd := &Command{
		Usage: "Commands for querying information about services and modifying them",
		Command: &cli.Command{
			Name: "services",
			Subcommands: []*cli.Command{
				createCmd.Package(),
				removeCmd.Package(),
				listCmd.Package(),
			},
		},
	}

	commands = append(commands, servicesCmd.Package())
}

func servicesCreateCmd(c *cli.Context) error {
	args := Arguments(c.Context)
	cfgSession := Session(c.Context)
	typeStr := c.String("type")

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	typ, err := primitives.NewServiceType(typeStr)
	if err != nil {
		return err
	}

	email := args["identifier"].(primitives.Email)

	clusterURL, err := cluster.GetURL()
	if err != nil {
		return err
	}

	apiToken, err := auth.NewAPIToken(email, clusterURL)
	if err != nil {
		return err
	}

	creds, err := apiToken.Credentials()
	if err != nil {
		return err
	}

	service, err := primitives.NewService(email, typ, creds.Package())
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	_, err = client.CreateService(c.Context, service)
	if err != nil {
		return err
	}

	tokenStr, err := apiToken.Marshal()
	if err != nil {
		return err
	}

	fmt.Printf("The service '%s' with type '%s' has been created. The following token "+
		"can be used to authenticate as that service:\n", service.Email, service.Type)

	fmt.Printf("Token: %s\n", tokenStr)
	fmt.Println("Remember: Please keep the token safe and share it securely as\n" +
		"it enables anyone with access to the cluster the ability to authenticate as this service.")

	return nil
}

func servicesRemoveCmd(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	ui := UI(c.Context)

	args := Arguments(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	email := args["identifier"].(primitives.Email)
	email.SetType(primitives.ServiceEmail)

	if !skipConfirm {
		err := ui.Confirm(fmt.Sprintf("Do you really want to delete %s and all of its tokens?", email))
		if err != nil {
			return err
		}
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	service, err := client.GetServiceByEmail(c.Context, email)
	if err != nil {
		return err
	}

	err = client.DeleteService(c.Context, service.ID)
	if err != nil {
		return err
	}

	fmt.Printf("The service '%s' has been deleted.", email)

	return nil
}

func servicesListCmd(c *cli.Context) error {
	ui := UI(c.Context)

	cfgSession := Session(c.Context)
	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	services, err := client.ListServices(c.Context)
	if err != nil {
		return err
	}

	header := []string{"Label", "Type", "Roles"}
	body := make([][]string, len(services))
	for i, s := range services {
		roleLabels := make([]string, len(s.Roles))
		for i, role := range s.Roles {
			roleLabels[i] = role.Label.String()
		}
		roles := strings.Join(roleLabels, ",")
		body[i] = []string{s.Email.String(), s.Type.String(), roles}
	}

	return ui.Table(header, body)
}
