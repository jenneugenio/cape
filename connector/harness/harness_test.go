// +build integration

package harness

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestHarness(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Can start the connector", func(t *testing.T) {
		ctx := context.Background()
		s, err := NewStack(ctx)
		gm.Expect(err).To(gm.BeNil())

		defer s.Teardown(ctx)

		url, err := s.ConnHarness.URL()
		gm.Expect(err).To(gm.BeNil())

		// ensure the connector is running
		resp, err := s.ConnHarness.server.Client().Get(url.String() + "/_healthz")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))
	})

	t.Run("Can start and stop the connector", func(t *testing.T) {
		ctx := context.Background()
		s, err := NewStack(ctx)
		gm.Expect(err).To(gm.BeNil())

		u, err := s.ConnHarness.URL()
		gm.Expect(err).To(gm.BeNil())

		// ensure the connector is running
		resp, err := s.ConnHarness.server.Client().Get(u.String() + "/_healthz")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))

		err = s.Teardown(ctx)
		gm.Expect(err).To(gm.BeNil())

		// httptest server is gone
		gm.Expect(s.ConnHarness.server).To(gm.BeNil())
	})
}
