package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
)

func TestTokens(t *testing.T) {
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

	t.Run("Create a token", func(t *testing.T) {
		gm.RegisterTestingT(t)

		token, err := client.CreateToken(ctx, m.Admin.User)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(token).ToNot(gm.BeNil())
	})

	t.Run("Can login with a token", func(t *testing.T) {
		gm.RegisterTestingT(t)

		token, err := client.CreateToken(ctx, m.Admin.User)
		gm.Expect(err).To(gm.BeNil())

		err = client.Logout(ctx, m.Admin.Token)
		gm.Expect(err).To(gm.BeNil())

		_, err = client.TokenLogin(ctx, token)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can remove a token", func(t *testing.T) {
		gm.RegisterTestingT(t)

		token, err := client.CreateToken(ctx, m.Admin.User)
		gm.Expect(err).To(gm.BeNil())

		err = client.RemoveToken(ctx, token.TokenID)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("Can list tokens", func(t *testing.T) {
		gm.RegisterTestingT(t)

		for i := 0; i < 10; i++ {
			_, err := client.CreateToken(ctx, m.Admin.User)
			gm.Expect(err).To(gm.BeNil())
		}

		tokenIDS, err := client.ListTokens(ctx, nil)
		gm.Expect(err).To(gm.BeNil())

		// TODO -- there are an extra 2 from the tests above
		// Would be nice to have a reset API, we could put this
		// into its own test block but that involves creating
		// and entirely new stack
		gm.Expect(len(tokenIDS)).To(gm.Equal(12))
	})
}
