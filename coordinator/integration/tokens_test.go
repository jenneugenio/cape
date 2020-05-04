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
}
