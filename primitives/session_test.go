package primitives

import (
	"context"
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
)

func TestNewSession(t *testing.T) {
	gm.RegisterTestingT(t)
	email, err := NewEmail("bob@bob.com")
	gm.Expect(err).To(gm.BeNil())

	pub, _, _ := ed25519.GenerateKey(nil)
	pkey := base64.New(pub)
	salt := base64.New([]byte("SALTSALTSALTSALT"))

	creds, err := NewCredentials(pkey, salt)
	gm.Expect(err).To(gm.BeNil())

	user, err := NewUser("bob", email, creds)
	gm.Expect(err).To(gm.BeNil())

	ti := time.Now().UTC().Add(time.Minute * 5)
	token := base64.New([]byte("random-string"))

	t.Run("new session", func(t *testing.T) {
		session, err := NewSession(user, Login)
		gm.Expect(err).To(gm.BeNil())
		session.SetToken(token, ti)

		gm.Expect(session.GetType()).To(gm.Equal(SessionType))
		gm.Expect(session.Credentials).ToNot(gm.BeNil())
		gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
		gm.Expect(session.Token).To(gm.Equal(token))
		gm.Expect(session.IdentityID).To(gm.Equal(user.ID))

		session, err = NewSession(user, Authenticated)
		gm.Expect(err).To(gm.BeNil())

		session.SetToken(token, ti)

		gm.Expect(session.GetType()).To(gm.Equal(SessionType))
		gm.Expect(session.Credentials).To(gm.BeNil())
		gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
		gm.Expect(session.Token).To(gm.Equal(token))
		gm.Expect(session.IdentityID).To(gm.Equal(user.ID))
	})

	t.Run("test encrypt decrytp", func(t *testing.T) {
		session, err := NewSession(user, Login)
		gm.Expect(err).To(gm.BeNil())

		session.SetToken(token, ti)

		key, err := crypto.NewBase64KeyURL(nil)
		gm.Expect(err).To(gm.BeNil())

		kms, err := crypto.LoadKMS(key)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := crypto.NewSecretBoxCodec(kms)
		gm.Expect(err).To(gm.BeNil())

		ctx := context.Background()
		by, err := session.Encrypt(ctx, codec)
		gm.Expect(err).To(gm.BeNil())

		newSession := &Session{}
		err = newSession.Decrypt(ctx, codec, by)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(newSession).To(gm.Equal(session))
	})
}
