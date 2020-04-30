// +build integration

package integration

import (
	"context"
	"github.com/manifoldco/go-base64"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/connector/harness"
)

func TestLogin(t *testing.T) {
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

	t.Run("can submit query that logs in", func(t *testing.T) {
		stream, err := connClient.Query(ctx, s.Manager.TestSource.Label, "SELECT * FROM transactions", 50, 0)
		gm.Expect(err).To(gm.BeNil())

		defer stream.Close()

		// NextRecord actually triggers the login
		stream.NextRecord()

		err = stream.Error()
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can still submit query that logs in", func(t *testing.T) {
		stream, err := connClient.Query(ctx, s.Manager.TestSource.Label, "SELECT * FROM transactions", 50, 0)
		gm.Expect(err).To(gm.BeNil())

		defer stream.Close()

		// NextRecord actually triggers the login
		stream.NextRecord()

		err = stream.Error()
		gm.Expect(err).To(gm.BeNil())
	})
}

func TestLoginFail(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	s, err := harness.NewStack(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer s.Teardown(ctx)

	connClient, err := s.ConnHarness.Client(base64.New([]byte("abcdefgh")))
	gm.Expect(err).To(gm.BeNil())

	defer connClient.Close()

	stream, err := connClient.Query(ctx, "test-datasource", "SELECT * FROM ALLDATA;", 50, 0)
	gm.Expect(err).To(gm.BeNil())

	// NextRecord actually triggers the login
	ok := stream.NextRecord()
	gm.Expect(ok).To(gm.BeFalse())
}
