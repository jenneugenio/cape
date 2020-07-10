// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/harness"
	"github.com/capeprivacy/cape/models"
	"github.com/capeprivacy/cape/primitives"
)

func TestRecoveries(t *testing.T) {
	gm.RegisterTestingT(t)

	ctx := context.Background()
	cfg, err := harness.NewConfig()
	gm.Expect(err).To(gm.BeNil())

	h, err := harness.NewHarness(cfg)
	gm.Expect(err).To(gm.BeNil())

	err = h.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	defer h.Teardown(ctx) // nolint: errcheck

	m := h.Manager()
	client, err := m.Setup(ctx)
	gm.Expect(err).To(gm.BeNil())

	unknownEmail := models.Email("unknwon@email.us")

	email := models.Email("jerry@jerry.berry")
	name := models.Name("Jerry Berry")

	password, err := primitives.NewPassword("hellotestingthisout")
	gm.Expect(err).To(gm.BeNil())

	user, _, err := client.CreateUser(ctx, name, email)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can recover account successfully", func(t *testing.T) {
		err := client.CreateRecovery(ctx, email)
		gm.Expect(err).To(gm.BeNil())

		mail := h.Mails()
		gm.Expect(len(mail)).To(gm.Equal(1))

		recovery := mail[0].Arguments["recovery"].(*primitives.Recovery)
		secret := mail[0].Arguments["secret"].(primitives.Password)

		err = client.AttemptRecovery(ctx, recovery.ID.String(), secret, password)
		gm.Expect(err).To(gm.BeNil())

		userClient, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		session, err := userClient.EmailLogin(ctx, email, password)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(session).ToNot(gm.BeNil())
	})

	t.Run("can request account recovery for unknown email", func(t *testing.T) {
		err := client.CreateRecovery(ctx, unknownEmail)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(h.Mails())).To(gm.Equal(1))
	})

	t.Run("can't recover with wrong id", func(t *testing.T) {
		err := client.CreateRecovery(ctx, email)
		gm.Expect(err).To(gm.BeNil())

		mail := h.Mails()
		gm.Expect(len(mail)).To(gm.Equal(2))

		secret := mail[1].Arguments["secret"].(primitives.Password)

		err = client.AttemptRecovery(ctx, user.ID, secret, password)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("can't recover account with wrong secret", func(t *testing.T) {
		err := client.CreateRecovery(ctx, email)
		gm.Expect(err).To(gm.BeNil())

		mail := h.Mails()
		gm.Expect(len(mail)).To(gm.Equal(3))

		recovery := mail[2].Arguments["recovery"].(*primitives.Recovery)

		err = client.AttemptRecovery(ctx, recovery.ID.String(), password, password)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
