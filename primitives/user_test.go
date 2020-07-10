package primitives

import (
	"context"
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database/crypto"
	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestUser(t *testing.T) {
	gm.RegisterTestingT(t)

	name, err := NewName("my-name")
	gm.Expect(err).To(gm.BeNil())

	creds := GenerateCredentials()

	email, err := NewEmail("email@email.com")
	gm.Expect(err).To(gm.BeNil())

	t.Run("new user", func(t *testing.T) {
		user, err := NewUser(name, email, creds)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(user.Email).To(gm.Equal(email))
		gm.Expect(user.Name).To(gm.Equal(name))
		gm.Expect(user.Credentials.Alg).To(gm.Equal(SHA256))
		gm.Expect(user.Credentials.Salt).To(gm.Equal(creds.Salt))
		gm.Expect(user.Credentials.Secret).To(gm.Equal(creds.Secret))
	})

	t.Run("test encrypt & decrypt", func(t *testing.T) {
		user, err := NewUser(name, email, creds)
		gm.Expect(err).To(gm.BeNil())

		key, err := crypto.NewBase64KeyURL(nil)
		gm.Expect(err).To(gm.BeNil())

		kms, err := crypto.LoadKMS(key)
		gm.Expect(err).To(gm.BeNil())

		defer kms.Close()

		codec := crypto.NewSecretBoxCodec(kms)
		ctx := context.Background()

		out, err := user.Encrypt(ctx, codec)
		gm.Expect(err).To(gm.BeNil())

		newUser := &User{}
		err = newUser.Decrypt(ctx, codec, out)
		gm.Expect(err).To(gm.BeNil())
		gm.Expect(newUser.Validate()).To(gm.BeNil())
		gm.Expect(newUser).To(gm.Equal(user))
	})

	badEmail := Email{
		Email: "this is not an email",
	}

	tests := []struct {
		testName string
		name     Name
		email    Email
		creds    *Credentials
	}{
		{
			"invalid name",
			Name("1"),
			email,
			creds,
		},
		{
			"invalid email",
			Name("HEYDHD"),
			badEmail,
			creds,
		},
		{
			"invalid credentials",
			Name("HEYDHD"),
			email,
			&Credentials{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := NewUser(tc.name, tc.email, tc.creds)
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(errors.CausedBy(err, InvalidUserCause)).To(gm.BeTrue())
		})
	}
}
