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
		Usage: "Pull data!",
		Command: &cli.Command{
			Name:   "pull",
			Action: pullDataCmd,
		},
	}

	commands = append(commands, pullCmd.Package())
}

func pullDataCmd(c *cli.Context) error {
	url, err := primitives.NewURL("https://localhost:8081")
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

	client, err := connector.NewClient(url, certPool)
	if err != nil {
		return err
	}

	source, err := primitives.NewLabel("source")
	if err != nil {
		return err
	}

	err = client.Query(c.Context, source, "SELECT * FROM ALLDATA;")
	if err != nil {
		return err
	}

	return nil
}
