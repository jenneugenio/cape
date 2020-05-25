// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

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

	email, err := primitives.NewEmail("newuser@email.com")
	gm.Expect(err).To(gm.BeNil())

	name, err := primitives.NewName("hello")
	gm.Expect(err).To(gm.BeNil())

	user, password, err := client.CreateUser(ctx, name, email)
	gm.Expect(err).To(gm.BeNil())

	userClient, err := h.Client()
	gm.Expect(err).To(gm.BeNil())

	_, err = userClient.EmailLogin(ctx, email, password)
	gm.Expect(err).To(gm.BeNil())

	t.Run("tokens for the admin", func(t *testing.T) {
		t.Run("Create a token", func(t *testing.T) {
			gm.RegisterTestingT(t)

			apiToken, token, err := client.CreateToken(ctx, m.Admin.User)
			gm.Expect(err).To(gm.BeNil())
			gm.Expect(token).ToNot(gm.BeNil())
			gm.Expect(apiToken).ToNot(gm.BeNil())

			// Credentials should never be leaked from the server
			creds, err := token.GetCredentials()
			gm.Expect(err).To(gm.BeNil())
			gm.Expect(creds).To(gm.BeNil())
			gm.Expect(token.Credentials).To(gm.BeNil())
		})

		t.Run("Can login with a token", func(t *testing.T) {
			gm.RegisterTestingT(t)

			apiToken, _, err := client.CreateToken(ctx, m.Admin.User)
			gm.Expect(err).To(gm.BeNil())

			err = client.Logout(ctx, m.Admin.Token)
			gm.Expect(err).To(gm.BeNil())

			_, err = client.TokenLogin(ctx, apiToken)
			gm.Expect(err).To(gm.BeNil())
		})

		t.Run("Can remove a token", func(t *testing.T) {
			gm.RegisterTestingT(t)

			apiToken, _, err := client.CreateToken(ctx, m.Admin.User)
			gm.Expect(err).To(gm.BeNil())

			err = client.RemoveToken(ctx, apiToken.TokenID)
			gm.Expect(err).To(gm.BeNil())
		})

		t.Run("Can list tokens", func(t *testing.T) {
			gm.RegisterTestingT(t)

			for i := 0; i < 10; i++ {
				_, _, err := client.CreateToken(ctx, m.Admin.User)
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
	})

	t.Run("tokens for another user", func(t *testing.T) {
		t.Run("can't list tokens you don't own", func(t *testing.T) {
			_, err = userClient.ListTokens(ctx, m.Admin.User)
			gm.Expect(err).NotTo(gm.BeNil())
		})

		t.Run("can't remove token you don't own", func(t *testing.T) {
			apiToken, _, err := client.CreateToken(ctx, m.Admin.User)
			gm.Expect(err).To(gm.BeNil())

			err = userClient.RemoveToken(ctx, apiToken.TokenID)
			gm.Expect(err).NotTo(gm.BeNil())
		})

		t.Run("admin can remove your token though", func(t *testing.T) {
			apiToken, _, err := userClient.CreateToken(ctx, user)
			gm.Expect(err).To(gm.BeNil())

			err = client.RemoveToken(ctx, apiToken.TokenID)
			gm.Expect(err).To(gm.BeNil())
		})

		t.Run("admin can list your tokens", func(t *testing.T) {
			_, _, err := userClient.CreateToken(ctx, user)
			gm.Expect(err).To(gm.BeNil())

			tokenIDS, err := client.ListTokens(ctx, user)
			gm.Expect(err).To(gm.BeNil())

			gm.Expect(len(tokenIDS)).To(gm.Equal(1))
		})
	})

	t.Run("service tokens", func(t *testing.T) {
		email, err := primitives.NewEmail("service:service@service.com")
		gm.Expect(err).To(gm.BeNil())

		url, err := primitives.NewURL("https://localhost:8081")
		gm.Expect(err).To(gm.BeNil())

		service, err := primitives.NewService(email, primitives.DataConnectorServiceType, url)
		gm.Expect(err).To(gm.BeNil())

		service, err = client.CreateService(ctx, service)
		gm.Expect(err).To(gm.BeNil())

		t.Run("admin can create service token", func(t *testing.T) {
			apiToken, token, err := client.CreateToken(ctx, service)
			gm.Expect(err).To(gm.BeNil())
			gm.Expect(token).ToNot(gm.BeNil())
			gm.Expect(apiToken).ToNot(gm.BeNil())
		})

		t.Run("regular use cannot create a service token", func(t *testing.T) {
			_, _, err := userClient.CreateToken(ctx, service)
			gm.Expect(err).NotTo(gm.BeNil())
		})
	})
}
