// +build integration

package harness

import (
	"context"
	"net/http"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestHarness(t *testing.T) {
	gm.RegisterTestingT(t)
	cfg, err := NewConfig()
	gm.Expect(err).To(gm.BeNil())

	t.Run("Can start the coordinator", func(t *testing.T) {
		ctx := context.Background()
		h, err := NewHarness(cfg)
		gm.Expect(err).To(gm.BeNil())
		defer h.Teardown(ctx) //nolint: errcheck

		err = h.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		url, err := h.URL()
		gm.Expect(err).To(gm.BeNil())

		// ensure the coordinator is running
		resp, err := http.Get(url.String() + "/_healthz")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))
	})

	t.Run("Can start and stop the coordinator", func(t *testing.T) {
		ctx := context.Background()
		h, err := NewHarness(cfg)
		gm.Expect(err).To(gm.BeNil())

		err = h.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		url, err := h.URL()
		gm.Expect(err).To(gm.BeNil())

		// ensure the coordinator is running
		resp, err := http.Get(url.String() + "/_healthz")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))

		err = h.Teardown(ctx)
		gm.Expect(err).To(gm.BeNil())

		// now this should fail
		_, err = http.Get(url.String() + "/_healthz")
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
