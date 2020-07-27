// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/models"
)

func TestRoles(t *testing.T) {
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
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	t.Run("Test getting your own role", func(t *testing.T) {
		r, err := client.MyRole(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(r.Label).To(gm.Equal(models.AdminRole))
	})

	t.Run("Test getting your own role in a project", func(t *testing.T) {
		_, err := client.CreateProject(ctx, "My Project", nil, "Who cares")
		gm.Expect(err).To(gm.BeNil())

		r, err := client.MyProjectRole(ctx, "my-project")
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(r.Label).To(gm.Equal(models.ProjectOwnerRole))
	})
}