// +build integration

package integration

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/harness"
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

	unknownEmail, err := primitives.NewEmail("unknwon@email.us")
	gm.Expect(err).To(gm.BeNil())

	email, err := primitives.NewEmail("jerry@jerry.berry")
	gm.Expect(err).To(gm.BeNil())

	name, err := primitives.NewName("Jerry Berry")
	gm.Expect(err).To(gm.BeNil())

	password, err := primitives.NewPassword("hellotestingthisout")
	gm.Expect(err).To(gm.BeNil())

	user, _, err := client.CreateUser(ctx, name, email)
	gm.Expect(err).To(gm.BeNil())

	workerToken, err := m.CreateWorker(ctx)
	gm.Expect(err).To(gm.BeNil())

	t.Run("can recover account successfully", func(t *testing.T) {
		err := client.CreateRecovery(ctx, email)
		gm.Expect(err).To(gm.BeNil())

		mail := h.Mails()
		gm.Expect(len(mail)).To(gm.Equal(1))

		recovery := mail[0].Arguments["recovery"].(*primitives.Recovery)
		secret := mail[0].Arguments["secret"].(primitives.Password)

		err = client.AttemptRecovery(ctx, recovery.ID, secret, password)
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

	t.Run("can't supply invalid email", func(t *testing.T) {
		err := client.CreateRecovery(ctx, primitives.Email{
			Email: "sdfsdf",
			Type:  primitives.UserEmail,
		})
		gm.Expect(err).ToNot(gm.BeNil())
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

		err = client.AttemptRecovery(ctx, recovery.ID, password, password)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("non-worker can't list recoveries", func(t *testing.T) {
		_, err := client.Recoveries(ctx)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("non-worker can't delete recoveries", func(t *testing.T) {
		ids := []database.ID{}
		for _, mail := range h.Mails() {
			recovery := mail.Arguments["recovery"].(*primitives.Recovery)
			ids = append(ids, recovery.ID)
		}

		err := client.DeleteRecoveries(ctx, ids)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("a worker can list recoveries", func(t *testing.T) {
		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		_, err = client.TokenLogin(ctx, workerToken)
		gm.Expect(err).To(gm.BeNil())

		recoveries, err := client.Recoveries(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(recoveries)).To(gm.Equal(2))
	})

	t.Run("a worker can delete recoveries", func(t *testing.T) {
		client, err := h.Client()
		gm.Expect(err).To(gm.BeNil())

		_, err = client.TokenLogin(ctx, workerToken)
		gm.Expect(err).To(gm.BeNil())

		recoveries, err := client.Recoveries(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(recoveries)).To(gm.Equal(2))

		ids := []database.ID{}
		for _, recovery := range recoveries {
			ids = append(ids, recovery.ID)
		}

		err = client.DeleteRecoveries(ctx, ids)
		gm.Expect(err).To(gm.BeNil())

		recoveries, err = client.Recoveries(ctx)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(len(recoveries)).To(gm.Equal(0))
	})
}
