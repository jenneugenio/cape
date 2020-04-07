package main

import (
	"bufio"
	"crypto/x509"
	"encoding/csv"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/dropoutlabs/cape/connector/client"
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
			Flags: []cli.Flag{
				outFlag(),
			},
		},
	}

	commands = append(commands, pullCmd.Package())
}

func pullDataCmd(c *cli.Context) error {
	cfgSession := Session(c.Context)
	args := Arguments(c.Context)

	outFile := c.String("out")

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	sourceLabel := args["source"].(primitives.Label)
	query := args["query"].(string)

	ctrlClient, err := cluster.Client()
	if err != nil {
		return err
	}

	token, err := cluster.Token()
	if err != nil {
		return err
	}

	source, err := ctrlClient.GetSourceByLabel(c.Context, sourceLabel)
	if err != nil {
		return err
	}

	if source.ServiceID == nil {
		return errors.New(NoLinkedService, "Source has not been linked to a data-connector")
	}

	service, err := ctrlClient.GetService(c.Context, *source.ServiceID)
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

	connClient, err := client.NewClient(service.Endpoint, token, certPool)
	if err != nil {
		return err
	}

	stream, err := connClient.Query(c.Context, sourceLabel, query)
	if err != nil {
		return err
	}

	file := os.Stdout
	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return err
		}
		file = f
	}

	w := bufio.NewWriter(file)
	writer := csv.NewWriter(w)

	for stream.NextRecord() {
		record := stream.Record()
		strs := record.ToStrings()

		err = writer.Write(strs)
		if err != nil {
			return err
		}

		writer.Flush()
		err = writer.Error()
		if err != nil {
			return err
		}
	}

	return stream.Error()
}
