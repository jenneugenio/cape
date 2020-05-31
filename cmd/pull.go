package main

import (
	"bufio"
	"crypto/x509"
	"encoding/csv"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/capeprivacy/cape/cmd/ui"
	"github.com/capeprivacy/cape/connector/client"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
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
		Arguments: []*Argument{SourceLabelArg, PullQueryArgument},
		Command: &cli.Command{
			Name:   "pull",
			Action: handleSessionOverrides(pullDataCmd),
			Flags: []cli.Flag{
				outFlag(),
				limitFlag(),
				offsetFlag(),
			},
		},
	}

	commands = append(commands, pullCmd.Package())
}

func pullDataCmd(c *cli.Context) error {
	cfgSession := Session(c.Context)
	outFile := c.String("out")
	limit := c.Int64("limit")
	offset := c.Int64("offset")

	cluster, err := cfgSession.Cluster()
	if err != nil {
		return err
	}

	sourceLabel := Arguments(c.Context, SourceLabelArg).(primitives.Label)
	query := Arguments(c.Context, PullQueryArgument).(string)

	provider := GetProvider(c.Context)
	coordClient, err := provider.Client(c.Context)
	if err != nil {
		return err
	}

	token, err := cluster.Token()
	if err != nil {
		return err
	}

	source, err := coordClient.GetSourceByLabel(c.Context, sourceLabel, nil)
	if err != nil {
		return err
	}

	if source.Service == nil {
		return errors.New(NoLinkedService, "Source has not been linked to a data-connector")
	}

	service, err := coordClient.GetService(c.Context, source.Service.ID)
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

	stream, err := connClient.Query(c.Context, sourceLabel, query, limit, offset)
	if err != nil {
		return err
	}

	if outFile != "" {
		f, err := os.Create(outFile)
		if err != nil {
			return err
		}
		file := f

		w := bufio.NewWriter(file)
		writer := csv.NewWriter(w)

		first := true

		for stream.NextRecord() {
			record := stream.Record()
			strs, err := record.ToStrings()
			if err != nil {
				return err
			}

			if first {
				schema := stream.Schema()
				var fieldNames []string
				for _, f := range schema.Fields {
					fieldNames = append(fieldNames, f.Name)
				}

				err = writer.Write(fieldNames)
				if err != nil {
					return err
				}

				writer.Flush()
				err = writer.Error()
				if err != nil {
					return err
				}

				first = false
			}

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

	// otherwise, print it in a nice table
	var body ui.TableBody
	for stream.NextRecord() {
		record := stream.Record()
		strs, err := record.ToStrings()
		if err != nil {
			return err
		}

		body = append(body, strs)
	}
	err = stream.Error()
	if err != nil {
		return err
	}

	u := provider.UI(c.Context)

	// We make the header second as the schema will not be on the stream until after we have received our first record
	schema := stream.Schema()
	var header ui.TableHeader
	if schema != nil && len(body) > 0 {
		for _, f := range schema.Fields {
			header = append(header, f.Name)
		}
		err := u.Table(header, body)
		if err != nil {
			return err
		}
	}

	return u.Template("\nFound {{ . | toString | faded }} record{{ . | pluralize \"s\"}}\n", len(body))
}
