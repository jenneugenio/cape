package integration

import (
	"context"
	"fmt"
	"github.com/dropoutlabs/cape/auth"
	"github.com/dropoutlabs/cape/controller"
	"github.com/dropoutlabs/cape/primitives"
	gm "github.com/onsi/gomega"
	"testing"
)

func TestSetup(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()

	tc, err := controller.NewTestController()
	gm.Expect(err).To(gm.BeNil())

	_, err = tc.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	fmt.Println("DB URL", tc.URL())
	defer tc.Teardown(ctx) // nolint: errcheck

	t.Run("Setup cape", func(t *testing.T) {
		gm.RegisterTestingT(t)
		client := controller.NewClient(tc.URL(), nil)
		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("ben", "ben@capeprivacy.com", creds.Package())
		gm.Expect(err).To(gm.BeNil())

		admin, err := client.Setup(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(admin.Name).To(gm.Equal(user.Name))
		gm.Expect(admin.Email).To(gm.Equal(user.Email))
	})

	t.Run("Setup cannot be called a second time", func(t *testing.T) {
		gm.RegisterTestingT(t)
		client := controller.NewClient(tc.URL(), nil)
		creds, err := auth.NewCredentials([]byte("jerryberrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("ben", "ben@capeprivacy.com", creds.Package())
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, user)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("Setup cannot be called a second time with different information", func(t *testing.T) {
		gm.RegisterTestingT(t)
		client := controller.NewClient(tc.URL(), nil)
		creds, err := auth.NewCredentials([]byte("berryjerrybuddyboy"), nil)
		gm.Expect(err).To(gm.BeNil())

		user, err := primitives.NewUser("justin", "justin@capeprivacy.com", creds.Package())
		gm.Expect(err).To(gm.BeNil())

		_, err = client.Setup(ctx, user)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
