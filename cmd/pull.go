package main

import (
	"crypto/x509"
	"io/ioutil"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/connector"
	errors "github.com/dropoutlabs/cape/partyerrors"
	"github.com/dropoutlabs/cape/primitives"
)

func init() {
	pullCmd := &Command{
		Usage: "The command pulls data from the given data source using the supplied query.",
		Examples: []*Example{
			{
				Example: "cape pull creditcards 'SELECT * FROM transactions'",
				Description: "Queries the table transactions from the data source creditcards. Policy " +
					"is applied to this query and some fields may be redacted/obfuscated or hidden in some " +
					"other privacy preserving manner.",
			},
			{
				Example: "cape pull creditcards transactions",
				Description: "Alias for querying data like 'SELECT * FROM transactions'. Any data that has " +
					"policy attached will be redacted/obfuscated or hidden in some other privacy preserving manner.",
			},
		},
		Arguments: []*Argument{LabelArg("source"), PullQueryArgument},
		Command: &cli.Command{
			Name:   "pull",
			Action: handleSessionOverrides(pullDataCmd),
		},
	}

	commands = append(commands, pullCmd.Package())
}

func pullDataCmd(c *cli.Context) error {
	cfgSession := Session(c.Context)
	args := Arguments(c.Context)

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	sourceLabel := args["source"].(primitives.Label)
	query := args["query"].(string)

	client, err := cluster.Client()
	if err != nil {
		return err
	}

	source, err := client.GetSourceByLabel(c.Context, sourceLabel)
	if err != nil {
		return err
	}

	service, err := client.GetService(c.Context, source.ServiceID)
	if err != nil {
		return err
	}

	cert, err := ioutil.ReadFile("connector/certs/localhost.crt")
	if err != nil {
		return err
	}

	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(cert)
	if !ok {
		return errors.New(BadCertificate, "Bad certificate for TLS")
	}

	connClient, err := connector.NewClient(service.Endpoint, certPool)
	if err != nil {
		return err
	}

	err = connClient.Query(c.Context, sourceLabel, query)
	if err != nil {
		return err
	}

	return nil
}
