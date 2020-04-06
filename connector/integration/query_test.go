// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	connHarness "github.com/dropoutlabs/cape/connector/harness"
	"github.com/dropoutlabs/cape/connector/sources"
	"github.com/dropoutlabs/cape/controller/harness"
)

func TestQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	_, err = m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	controllerURL, err := m.URL()
	gm.Expect(err).To(gm.BeNil())

	connCfg, err := connHarness.NewConfig(controllerURL)
	gm.Expect(err).To(gm.BeNil())

	connH, err := connHarness.NewHarness(connCfg)
	gm.Expect(err).To(gm.BeNil())

	err = connH.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer connH.Teardown(ctx) // nolint: errcheck

	connectorURL, err := connH.URL()
	gm.Expect(err).To(gm.BeNil())

	err = m.CreateService(ctx, connH.APIToken(), connectorURL)
	gm.Expect(err).To(gm.BeNil())

	err = m.CreateSource(ctx, connH.SourceCredentials(), m.Connector.ID)
	gm.Expect(err).To(gm.BeNil())

	err = m.CreatePolicy(ctx, "./testdata/policy.yaml")
	gm.Expect(err).To(gm.BeNil())

	connClient, err := connH.Client(m.Admin.Token)
	gm.Expect(err).To(gm.BeNil())

	defer connClient.Close()

	query := "SELECT * FROM transactions"
	stream, err := connClient.Query(context.Background(), m.TestSource.Label, query)
	gm.Expect(err).To(gm.BeNil())

	defer stream.Close()

	expectedRows, err := sources.GetExpectedRows(ctx, connH.SourceCredentials().ToURL(), query, nil)
	gm.Expect(err).To(gm.BeNil())

	i := 0
	for stream.NextRecord() {
		if i > 5 {
			break
		}
		err := stream.Error()
		gm.Expect(err).To(gm.BeNil())

		record := stream.Record()
		// could check row to row but this is easier to see
		// if there are any errors
		for j, val := range record.Values() {
			gm.Expect(val).To(gm.Equal(expectedRows[i][j]))
		}
		i++
	}
}
