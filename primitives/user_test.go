package primitives

import (
	"crypto/ed25519"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestUser(t *testing.T) {
	gm.RegisterTestingT(t)

	pub, _, _ := ed25519.GenerateKey(nil)
	pkey := base64.New(pub)
	salt := base64.New([]byte("SALTSALTSALTSALT"))

	creds, err := NewCredentials(pkey, salt)
	gm.Expect(err).To(gm.BeNil())

	email, err := NewEmail("email@email.com")
	gm.Expect(err).To(gm.BeNil())

	t.Run("new user", func(t *testing.T) {
		name := Name("my-name")
		user, err := NewUser(name, email, creds)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(user.Email).To(gm.Equal(email))
		gm.Expect(user.Name).To(gm.Equal(name))
		gm.Expect(user.Credentials.Alg).To(gm.Equal(EDDSA))
		gm.Expect(user.Credentials.Salt).To(gm.Equal(salt))
		gm.Expect(user.Credentials.PublicKey).To(gm.Equal(pkey))
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
