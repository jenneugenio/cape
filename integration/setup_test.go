// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller/harness"
	"github.com/dropoutlabs/cape/primitives"
)

func TestSetup(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err)

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	t.Run("Setup cape", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("ben@capeprivacy.com")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("ben", email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		admin, err := client.Setup(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(admin.Name).To(gm.Equal(user.Name))
		gm.Expect(admin.Email).To(gm.Equal(user.Email))
	})

	t.Run("Setup cannot be called a second time", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("ben@capeprivacy.com")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("ben", email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, user)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Setup cannot be called a second time with different information", func(t *testing.T) {
		gm.RegisterTestingT(t)

		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		creds, err := auth.NewCredentials([]byte("berryjerrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		email, err := primitives.NewEmail("justin@capeprivacy.com")
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("justin", email, creds.Package())
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, user)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
