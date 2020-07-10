package primitives

import (
	"context"
	"testing"
	"time"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestSession(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user, err := GenerateUser("bob", "test@email.com")
	gm.Expect(err).To(gm.BeNil())

	_, token, err := GenerateToken(user)
	gm.Expect(err).To(gm.BeNil())

	ti := time.Now().UTC().Add(time.Minute * 5)
	sessionToken := base64.New([]byte("random-string"))

	t.Run("validate", func(t *testing.T) {
		tests := []struct {
			name  string
			fn    func() (*Session, error)
			cause *errors.Cause
		}{
			{
				name: "valid session owned by user",
				fn: func() (*Session, error) {
					return NewSession(user)
				},
			},
			{
				name: "valid session owned by token",
				fn: func() (*Session, error) {
					return NewSession(token)
				},
			},
			{
				name: "invalid id",
				fn: func() (*Session, error) {
					session, err := NewSession(user)
					if err != nil {
						return nil, err
					}

					session.ID = database.EmptyID
					return session, nil
				},
				cause: &InvalidSessionCause,
			},
			{
				name: "invalid user id",
				fn: func() (*Session, error) {
					session, err := NewSession(user)
					if err != nil {
						return nil, err
					}

					session.UserID = ""
					return session, nil
				},
				cause: &InvalidSessionCause,
			},
			{
				name: "user id is not a user",
				fn: func() (*Session, error) {
					session, err := NewSession(user)
					if err != nil {
						return nil, err
					}

					session.UserID = ""
					return session, nil
				},
				cause: &InvalidSessionCause,
			},
			{
				name: "invalid owner id",
				fn: func() (*Session, error) {
					session, err := NewSession(user)
					if err != nil {
						return nil, err
					}

					session.OwnerID = ""
					return session, nil
				},
				cause: &InvalidSessionCause,
			},
			{
				name: "owner id is not a token or user",
				fn: func() (*Session, error) {
					session, err := NewSession(user)
					if err != nil {
						return nil, err
					}

					session.OwnerID = ""
					return session, nil
				},
				cause: &InvalidSessionCause,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				s, err := tc.fn()
				gm.Expect(err).To(gm.BeNil())

				err = s.Validate()
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(err).To(gm.BeNil())
			})
		}
	})

	t.Run("new session", func(t *testing.T) {
		session, err := NewSession(user)
		gm.Expect(err).To(gm.BeNil())
		session.SetToken(sessionToken, ti)

		gm.Expect(session.GetType()).To(gm.Equal(SessionType))
		gm.Expect(session.ExpiresAt).To(gm.Equal(ti))
		gm.Expect(session.Token).To(gm.Equal(sessionToken))
		gm.Expect(session.UserID).To(gm.Equal(user.ID.String()))
		gm.Expect(session.OwnerID).To(gm.Equal(user.ID.String()))
	})

	t.Run("test encrypt decrytp", func(t *testing.T) {
		session, err := NewSession(user)
		gm.Expect(err).To(gm.BeNil())

		session.SetToken(sessionToken, ti)

		key, err := crypto.NewBase64KeyURL(nil)
		gm.Expect(err).To(gm.BeNil())

		kms, err := crypto.LoadKMS(key)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := crypto.NewSecretBoxCodec(kms)

		ctx := context.Background()
		by, err := session.Encrypt(ctx, codec)
		gm.Expect(err).To(gm.BeNil())

		newSession := &Session{}
		err = newSession.Decrypt(ctx, codec, by)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(newSession).To(gm.Equal(session))
	})
}
