package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/connector/harness"
	"github.com/capeprivacy/cape/primitives"
	"github.com/capeprivacy/cape/version"
)

func TestVersion(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	url, err := primitives.NewURL("http://localhost:8080")
	gm.Expect(err).To(gm.BeNil())

	connCfg, err := harness.NewConfig(url)
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(connCfg)
	gm.Expect(err).To(gm.BeNil())

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) //nolint: errcheck

	client, err := h.Client(nil)
	gm.Expect(err).To(gm.BeNil())

	defer client.Close()

	res, err := client.Version(ctx)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(res.Version).To(gm.Equal(version.Version))
	gm.Expect(res.BuildDate).To(gm.Equal(version.BuildDate))
}
