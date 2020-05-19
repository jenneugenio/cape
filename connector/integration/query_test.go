// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/capeprivacy/cape/connector/harness"
	"github.com/capeprivacy/cape/connector/proto"
	"github.com/capeprivacy/cape/connector/sources"
	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
)

func TestQuery(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	s, err := harness.NewStack(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer s.Teardown(ctx)

	err = s.Manager.CreateSource(ctx, s.ConnHarness.SourceCredentials(), s.Manager.Connector.ID)
	gm.Expect(err).To(gm.BeNil())

	err = s.Manager.CreatePolicy(ctx, "./testdata/policy.yaml")
	gm.Expect(err).To(gm.BeNil())

	connClient, err := s.ConnHarness.Client(s.Manager.Admin.Token)
	gm.Expect(err).To(gm.BeNil())

	defer connClient.Close()

	t.Run("query without transform", func(t *testing.T) {
		query := "SELECT * FROM transactions"
		stream, err := connClient.Query(context.Background(), s.Manager.TestSource.Label, query, 50, 50)
		gm.Expect(err).To(gm.BeNil())

		defer stream.Close()

		expectedRows, err := sources.GetExpectedRows(ctx, s.ConnHarness.SourceCredentials().ToURL(), query+" LIMIT 50 OFFSET 50", nil)
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
	})

	err = s.Manager.CreatePolicy(ctx, "./testdata/policy_transform_value.yaml")
	gm.Expect(err).To(gm.BeNil())

	t.Run("query with transform", func(t *testing.T) {
		query := "SELECT * FROM transactions"
		stream, err := connClient.Query(context.Background(), s.Manager.TestSource.Label, query, 50, 50)
		gm.Expect(err).To(gm.BeNil())

		defer stream.Close()

		expectedRows, err := sources.GetExpectedRows(ctx, s.ConnHarness.SourceCredentials().ToURL(), query+" LIMIT 50 OFFSET 50", nil)
		gm.Expect(err).To(gm.BeNil())

		i := 0
		var schema *proto.Schema
		for stream.NextRecord() {
			if schema == nil {
				schema = stream.Schema()
			}

			record := stream.Record()
			// could check row to row but this is easier to see
			// if there are any errors
			for j, val := range record.Values() {
				if schema.Fields[j].Name == "value" {
					gm.Expect(val).To(gm.Equal(expectedRows[i][j].(float64) + 1))
				} else {
					gm.Expect(val).To(gm.Equal(expectedRows[i][j]))
				}
			}
			i++
		}
		gm.Expect(stream.Error()).To(gm.BeNil())
	})
}

func TestQueryDenied(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	s, err := harness.NewStack(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer s.Teardown(ctx)

	err = s.Manager.CreateSource(ctx, s.ConnHarness.SourceCredentials(), s.Manager.Connector.ID)
	gm.Expect(err).To(gm.BeNil())

	connClient, err := s.ConnHarness.Client(s.Manager.Admin.Token)
	gm.Expect(err).To(gm.BeNil())

	defer connClient.Close()

	query := "SELECT * FROM transactions"
	stream, err := connClient.Query(context.Background(), s.Manager.TestSource.Label, query, 50, 0)
	gm.Expect(err).To(gm.BeNil())

	defer stream.Close()

	ok := stream.NextRecord()
	gm.Expect(ok).To(gm.BeFalse())

	err = stream.Error()
	gm.Expect(err).ToNot(gm.BeNil())
	gm.Expect(err.(*errors.Error).Messages[0]).To(gm.Equal("No policies match the provided query"))
}
