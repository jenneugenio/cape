package primitives

import (
	"context"
	"testing"
	"time"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestRecovery(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user, err := GenerateUser("hi", "hi@hi.hi")
	gm.Expect(err).To(gm.BeNil())

	creds := GenerateCredentials()

	t.Run("Validate", func(t *testing.T) {
		tests := []struct {
			name  string
			fn    func() (*Recovery, error)
			cause *errors.Cause
		}{
			{
				name: "valid recovery",
				fn: func() (*Recovery, error) {
					return NewRecovery(user.ID.String(), creds)
				},
			},
			{
				name: "invalid id",
				fn: func() (*Recovery, error) {
					r, err := NewRecovery(user.ID.String(), creds)
					if err != nil {
						return nil, err
					}

					r.ID = database.EmptyID
					return r, nil
				},
				cause: &InvalidRecoveryCause,
			},
			{
				name: "invalid user id",
				fn: func() (*Recovery, error) {
					r, err := NewRecovery(user.ID.String(), creds)
					if err != nil {
						return nil, err
					}

					r.UserID = ""
					return r, nil
				},
				cause: &InvalidRecoveryCause,
			},
			{
				name: "missing credentials",
				fn: func() (*Recovery, error) {
					r, err := NewRecovery(user.ID.String(), creds)
					if err != nil {
						return nil, err
					}

					r.Credentials = nil
					return r, nil
				},
				cause: &InvalidRecoveryCause,
			},
			{
				name: "bad credentials",
				fn: func() (*Recovery, error) {
					badcreds := GenerateCredentials()
					gm.Expect(err).To(gm.BeNil())

					r, err := NewRecovery(user.ID.String(), badcreds)
					if err != nil {
						return nil, err
					}

					r.Credentials.Alg = CredentialsAlgType("sdfs")
					return r, nil
				},
				cause: &InvalidRecoveryCause,
			},
			{
				name: "expires at is zero value",
				fn: func() (*Recovery, error) {
					r, err := NewRecovery(user.ID.String(), creds)
					if err != nil {
						return nil, err
					}

					r.ExpiresAt = time.Time{}
					return r, nil
				},
				cause: &InvalidRecoveryCause,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				r, err := tc.fn()
				gm.Expect(err).To(gm.BeNil())

				err = r.Validate()
				if tc.cause != nil {
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(err).To(gm.BeNil())
			})
		}
	})

	t.Run("expired returns true if expire at exceeded", func(t *testing.T) {
		r, err := GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		r.ExpiresAt = time.Now().UTC().Add(-1 * time.Minute)
		gm.Expect(r.Expired()).To(gm.BeTrue())
	})

	t.Run("expired returns true if expire at exceeded", func(t *testing.T) {
		r, err := GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(r.Expired()).To(gm.BeFalse())
	})

	t.Run("encrypt & decrypt", func(t *testing.T) {
		creds := GenerateCredentials()

		r, err := NewRecovery(user.ID.String(), creds)
		gm.Expect(err).To(gm.BeNil())

		key, err := crypto.NewBase64KeyURL(nil)
		gm.Expect(err).To(gm.BeNil())

		kms, err := crypto.LoadKMS(key)
		gm.Expect(err).To(gm.BeNil())
		defer kms.Close()

		codec := crypto.NewSecretBoxCodec(kms)

		ct, err := r.Encrypt(context.TODO(), codec)
		gm.Expect(err).To(gm.BeNil())

		out := &Recovery{}
		err = out.Decrypt(context.TODO(), codec, ct)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(r).To(gm.Equal(out))
	})
}
