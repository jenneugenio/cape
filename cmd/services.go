package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/cmd/ui"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
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
				Example: "cape services create --type data-connector --endpoint connector.prod.mycompany.com service:dc@prod.mycompany.com",
				Description: "Creates a service called service:dc@prod.mycompany.com with the endpoint connector.prod.mycompany.com " +
					"representing a Cape data connector.",
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
				dataConnectorEndpointFlag(),
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
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	typeStr := c.String("type")
	endpointStr := c.String("endpoint")

	typ, err := primitives.NewServiceType(typeStr)
	if err != nil {
		return err
	}

	var endpoint *primitives.URL
	if typ == primitives.DataConnectorServiceType {
		url, err := primitives.NewURL(endpointStr)
		if err != nil {
			return errors.New(MustSupplyEndpoint, "Must supply a valid endpoint when creating a data-connector service")
		}
		endpoint = url
	}

	email := Arguments(c.Context, ServiceIdentifierArg).(primitives.Email)

	creds, err := auth.RandomCredentials()
	if err != nil {
		return err
	}

	pCreds, err := creds.Package()
	if err != nil {
		return err
	}

	// TODO -- services need to lose credentials
	service, err := primitives.NewService(email, typ, endpoint, pCreds)
	if err != nil {
		return err
	}

	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	service, err = client.CreateService(c.Context, service)
	if err != nil {
		return err
	}

	apiToken, err := client.NewToken(c.Context, service, creds)
	if err != nil {
		return err
	}

	tokenStr, err := apiToken.Marshal()
	if err != nil {
		return err
	}

	args := struct {
		Email string
		Type  string
	}{
		service.Email.String(),
		service.Type.String(),
	}

	err = u.Template("The service {{ .Email | bold }} with type {{ .Type | bold }} has been created. The following token "+
		"can be used to authenticate as that service:\n\n", args)
	if err != nil {
		return err
	}

	err = u.Details(ui.Details{
		"Token": tokenStr,
	})
	if err != nil {
		return err
	}

	return u.Notify(ui.Remember, "Please keep the token safe and share it only over secure channels.")
}

func servicesRemoveCmd(c *cli.Context) error {
	skipConfirm := c.Bool("yes")
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	email := Arguments(c.Context, ServiceIdentifierArg).(primitives.Email)
	email.SetType(primitives.ServiceEmail)

	if !skipConfirm {
		err := u.Confirm(fmt.Sprintf("Do you really want to delete %s and all of its tokens?", email))
		if err != nil {
			return err
		}
	}

	client, err := provider.Client(c.Context)
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

	return u.Template("The service {{ . | bold }} has been deleted.", email.String())
}

func servicesListCmd(c *cli.Context) error {
	provider := GetProvider(c.Context)
	u := provider.UI(c.Context)

	client, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	services, err := client.ListServices(c.Context)
	if err != nil {
		return err
	}

	if len(services) > 0 {
		header := []string{"Label", "Type", "Endpoint", "Roles"}
		body := make([][]string, len(services))
		for i, s := range services {
			roleLabels := make([]string, len(s.Roles))
			for i, role := range s.Roles {
				roleLabels[i] = role.Label.String()
			}
			roles := strings.Join(roleLabels, ",")

			endpoint := ""
			if s.Endpoint != nil {
				endpoint = s.Endpoint.String()
			}

			body[i] = []string{s.Email.String(), s.Type.String(), endpoint, roles}
		}

		err = u.Table(header, body)
		if err != nil {
			return err
		}
	}

	return u.Template("\nFound {{ . | toString | faded }} services{{ . | pluralize \"s\"}}\n", len(services))
}
