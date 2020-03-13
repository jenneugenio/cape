package controller

import (
	"context"
	gm "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestControllerLifecycle(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Can start the controller", func(t *testing.T) {
		ctx := context.Background()
		tc, err := NewTestController()
		gm.Expect(err).To(gm.BeNil())
		defer tc.Teardown(ctx) //nolint: errcheck

		_, err = tc.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		// ensure the controller is running
		resp, err := http.Get("http://localhost:8081/health")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))
	})

	t.Run("Can start and stop the controller", func(t *testing.T) {
		ctx := context.Background()
		tc, err := NewTestController()
		gm.Expect(err).To(gm.BeNil())

		_, err = tc.Setup(ctx)
		gm.Expect(err).To(gm.BeNil())

		// ensure the controller is running
		resp, err := http.Get("http://localhost:8081/health")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(resp.StatusCode).To(gm.Equal(200))

		err = tc.Teardown(ctx)
		gm.Expect(err).To(gm.BeNil())

		// now this should fail
		_, err = http.Get("http://localhost:8081/health")
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
