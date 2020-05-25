package primitives

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestToken(t *testing.T) {
	gm.RegisterTestingT(t)

	_, user, err := GenerateUser("hello", "bob@hello.com")
	gm.Expect(err).To(gm.BeNil())

	roleID, err := database.GenerateID(RoleType)
	gm.Expect(err).To(gm.BeNil())

	serviceID, err := database.GenerateID(ServicePrimitiveType)
	gm.Expect(err).To(gm.BeNil())

	t.Run("create token", func(t *testing.T) {
		tests := []struct {
			name       string
			identityID database.ID
			creds      *Credentials
			cause      *errors.Cause
		}{
			{
				name:       "valid parameters",
				identityID: user.ID,
				creds:      user.Credentials,
			},
			{
				name:       "invalid identity id",
				identityID: database.EmptyID,
				creds:      user.Credentials,
				cause:      &database.InvalidIDCause,
			},
			{
				name:       "missing credentials",
				identityID: user.ID,
				creds:      nil,
				cause:      &InvalidTokenCause,
			},
			{
				name:       "wrong id type",
				identityID: roleID,
				creds:      user.Credentials,
				cause:      &InvalidTokenCause,
			},
			{
				name:       "valid for service",
				identityID: serviceID,
				creds:      user.Credentials,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				token, err := NewToken(tc.identityID, tc.creds)
				if tc.cause != nil {
					gm.Expect(err).ToNot(gm.BeNil())
					gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
					return
				}

				gm.Expect(token.IdentityID).To(gm.Equal(tc.identityID))
				gm.Expect(token.Credentials.Alg).To(gm.Equal(tc.creds.Alg))
				gm.Expect(token.Credentials.Salt).To(gm.Equal(tc.creds.Salt))
				gm.Expect(token.Credentials.Secret).To(gm.Equal(tc.creds.Secret))
			})
		}
	})

	t.Run("can encrypt and decrypt", func(t *testing.T) {
		_, token, err := GenerateToken(user)
		gm.Expect(err).To(gm.BeNil())

		key, err := crypto.NewBase64KeyURL(nil)
		gm.Expect(err).To(gm.BeNil())

		kms, err := crypto.LoadKMS(key)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := crypto.NewSecretBoxCodec(kms)
		ctx := context.Background()

		out, err := token.Encrypt(ctx, codec)
		gm.Expect(err).To(gm.BeNil())

		newToken := &Token{}
		err = newToken.Decrypt(ctx, codec, out)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(newToken.Validate()).To(gm.BeNil())
		gm.Expect(newToken).To(gm.Equal(token))
	})
}
