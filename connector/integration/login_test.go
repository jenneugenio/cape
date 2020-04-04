// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	connHarness "github.com/dropoutlabs/cape/connector/harness"
	"github.com/dropoutlabs/cape/controller/harness"
	"github.com/dropoutlabs/cape/primitives"
)

func TestLogin(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx)

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

	defer connH.Teardown(ctx)

	connectorURL, err := connH.URL()
	gm.Expect(err).To(gm.BeNil())

	err = m.CreateService(ctx, connH.APIToken(), connectorURL)
	gm.Expect(err).To(gm.BeNil())

	err = m.CreateSource(ctx, connH.SourceCredentials(), m.Connector.ID)
	gm.Expect(err).To(gm.BeNil())

	connClient, err := connH.Client(m.Admin.Token)
	gm.Expect(err).To(gm.BeNil())

	defer connClient.Close()

	t.Run("can submit query that logs in", func(t *testing.T) {
		stream, err := connClient.Query(ctx, m.TestSource.Label, "SELECT * FROM transactions")
		gm.Expect(err).To(gm.BeNil())

		defer stream.Close()

		// NextRecord actually triggers the login
		ok := stream.NextRecord()
		gm.Expect(ok).To(gm.BeTrue())

		err = stream.Error()
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can still submit query that logs in", func(t *testing.T) {
		stream, err := connClient.Query(ctx, m.TestSource.Label, "SELECT * FROM transactions")
		gm.Expect(err).To(gm.BeNil())

		defer stream.Close()

		// NextRecord actually triggers the login
		ok := stream.NextRecord()
		gm.Expect(ok).To(gm.BeTrue())

		err = stream.Error()
		gm.Expect(err).To(gm.BeNil())
	})
}

func TestLoginFail(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	url, err := primitives.NewURL("http://localhost:8080")
	gm.Expect(err).To(gm.BeNil())

	cfg, err := connHarness.NewConfig(url)
	gm.Expect(err).To(gm.BeNil())

	h, err := connHarness.NewHarness(cfg)
	gm.Expect(err).To(gm.BeNil())

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx)

	connClient, err := h.Client(base64.New([]byte("abcdefgh")))
	gm.Expect(err).To(gm.BeNil())

	defer connClient.Close()

	stream, err := connClient.Query(ctx, "test-datasource", "SELECT * FROM ALLDATA;")
	gm.Expect(err).To(gm.BeNil())

	// NextRecord actually triggers the login
	ok := stream.NextRecord()
	gm.Expect(ok).To(gm.BeFalse())
}
