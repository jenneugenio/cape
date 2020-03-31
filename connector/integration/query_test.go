// +build integration

package integration

import (
	"testing"
	"context"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/controller/harness"
	connHarness "github.com/dropoutlabs/cape/connector/harness"
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

	connCfg := &connHarness.Config{Token: m.Connector.Token}
	connH, err := connHarness.NewHarness(connCfg)
	gm.Expect(err).To(gm.BeNil())

	err = connH.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	connClient, err := connH.Client()
	gm.Expect(err).To(gm.BeNil())

	err = connClient.Query(ctx, "test-datasource", "SELECT * FROM ALLDATA;")
	gm.Expect(err).To(gm.BeNil())
}