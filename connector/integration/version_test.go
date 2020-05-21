// +build integration

package integration

import (
	"context"
	"github.com/capeprivacy/cape/connector/harness"
	"github.com/capeprivacy/cape/version"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestVersion(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	s, err := harness.NewStack(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer s.Teardown(ctx)

	client, err := s.ConnHarness.Client(s.Manager.Admin.Token)
	gm.Expect(err).To(gm.BeNil())

	defer client.Close()

	res, err := client.Version(ctx)
	gm.Expect(err).To(gm.BeNil())
	gm.Expect(res.Version).To(gm.Equal(version.Version))
	gm.Expect(res.BuildDate).To(gm.Equal(version.BuildDate))
}
