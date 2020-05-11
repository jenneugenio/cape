// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/auth"
	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/primitives"
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

	email, _ := primitives.NewEmail("newuser@email.com")
	secret := []byte("verysecretsecret")

	creds, err := auth.NewCredentials(secret, nil)
	gm.Expect(err).To(gm.BeNil())

	pCreds, err := creds.Package()
	gm.Expect(err).To(gm.BeNil())

	user, err := primitives.NewUser("New User", email, pCreds)
	gm.Expect(err).To(gm.BeNil())

	user, err = client.CreateUser(ctx, user)
	gm.Expect(err).To(gm.BeNil())

	userClient, err := h.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = userClient.EmailLogin(ctx, email, secret)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can't list tokens you don't own", func(t *testing.T) {
		_, err = userClient.ListTokens(ctx, m.Admin.User)
		gm.Expect(err).NotTo(gm.BeNil())
	})

	t.Run("can't remove token you don't own", func(t *testing.T) {
		token, err := client.CreateToken(ctx, m.Admin.User)
		gm.Expect(err).To(gm.BeNil())

		err = userClient.RemoveToken(ctx, token.TokenID)
		gm.Expect(err).NotTo(gm.BeNil())
	})

	t.Run("admin can remove your token though", func(t *testing.T) {
		token, err := userClient.CreateToken(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		err = client.RemoveToken(ctx, token.TokenID)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("admin can list your tokens", func(t *testing.T) {
		_, err := userClient.CreateToken(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		tokenIDS, err := client.ListTokens(ctx, user)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(tokenIDS)).To(gm.Equal(1))
	})

	email, _ = primitives.NewEmail("service:service@service.com")
	url, _ := primitives.NewURL("https://localhost:8081")
	service, err := primitives.NewService(email, primitives.DataConnectorServiceType, url)
	gm.Expect(err).To(gm.BeNil())

	service, err = client.CreateService(ctx, service)
	gm.Expect(err).To(gm.BeNil())

	t.Run("admin can create service token", func(t *testing.T) {
		token, err := client.CreateToken(ctx, service)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(token).ToNot(gm.BeNil())
	})

	t.Run("regular use cannot create a service token", func(t *testing.T) {
		_, err := userClient.CreateToken(ctx, service)
		gm.Expect(err).NotTo(gm.BeNil())
	})
}
