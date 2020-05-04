// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	connHarness "github.com/capeprivacy/cape/connector/harness"
	"github.com/capeprivacy/cape/connector/sources"
	"github.com/capeprivacy/cape/coordinator/harness"
	errors "github.com/capeprivacy/cape/partyerrors"
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

	coordinatorURL, err := m.URL()
	gm.Expect(err).To(gm.BeNil())

	connCfg, err := connHarness.NewConfig(coordinatorURL)
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
	stream, err := connClient.Query(context.Background(), m.TestSource.Label, query, 50, 50)
	gm.Expect(err).To(gm.BeNil())

	defer stream.Close()

	expectedRows, err := sources.GetExpectedRows(ctx, connH.SourceCredentials().ToURL(), query+" LIMIT 50 OFFSET 50", nil)
	gm.Expect(err).To(gm.BeNil())

	i := 0
	for stream.NextRecord() {
		record := stream.Record()
		// could check row to row but this is easier to see
		// if there are any errors
		for j, val := range record.Values() {
			gm.Expect(val).To(gm.Equal(expectedRows[i][j]))
		}
		i++
	}
	gm.Expect(stream.Error()).To(gm.BeNil())
}

func TestQueryDenied(t *testing.T) {
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

	coordinatorURL, err := m.URL()
	gm.Expect(err).To(gm.BeNil())

	connCfg, err := connHarness.NewConfig(coordinatorURL)
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

	connClient, err := connH.Client(m.Admin.Token)
	gm.Expect(err).To(gm.BeNil())

	defer connClient.Close()

	query := "SELECT * FROM transactions"
	stream, err := connClient.Query(context.Background(), m.TestSource.Label, query, 50, 0)
	gm.Expect(err).To(gm.BeNil())

	defer stream.Close()

	ok := stream.NextRecord()
	gm.Expect(ok).To(gm.BeFalse())

	err = stream.Error()
	gm.Expect(err).ToNot(gm.BeNil())
	gm.Expect(err.(*errors.Error).Messages[0]).To(gm.Equal("No policies match the provided query"))
}
