// +build integration

package harness

import (
	"context"
	"testing"

	"github.com/dropoutlabs/cape/primitives"
	gm "github.com/onsi/gomega"
)

func TestHarness(t *testing.T) {
	gm.RegisterTestingT(t)

	controllerURL, err := primitives.NewURL("http://localhost:8080")
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can start the connector", func(t *testing.T) {
		ctx := context.Background()

		cfg := &Config{ControllerURL: controllerURL}

		h, err := NewHarness(cfg)
		gm.Expect(err).To(gm.BeNil())
		defer h.Teardown(ctx) //nolint: errcheck

		err = h.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		url, err := h.URL()
		gm.Expect(err).To(gm.BeNil())

		// ensure the connector is running
		resp, err := h.server.Client().Get(url.String() + "/_healthz")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))
	})

	t.Run("Can start and stop the connector", func(t *testing.T) {
		ctx := context.Background()

		cfg := &Config{ControllerURL: controllerURL}

		h, err := NewHarness(cfg)
		gm.Expect(err).To(gm.BeNil())

		err = h.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		url, err := h.URL()
		gm.Expect(err).To(gm.BeNil())

		// ensure the connector is running
		resp, err := h.server.Client().Get(url.String() + "/_healthz")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))

		err = h.Teardown(ctx)
		gm.Expect(err).To(gm.BeNil())

		// httptest server is gone
		gm.Expect(h.server).To(gm.BeNil())
	})
}
