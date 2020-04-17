package primitives

import (
	"crypto/ed25519"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestCredentials(t *testing.T) {
	gm.RegisterTestingT(t)

	pub, _, _ := ed25519.GenerateKey(nil)
	pkey := base64.New(pub)
	salt := base64.New([]byte("SALTSALTSALTSALT"))

	t.Run("create credentials", func(t *testing.T) {
		creds, err := NewCredentials(pkey, salt)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(creds.PublicKey).To(gm.Equal(pkey))
		gm.Expect(creds.Salt).To(gm.Equal(salt))
	})

	tests := []struct {
		name string
		pkey *base64.Value
		salt *base64.Value
	}{
		{
			"invalid private key",
			nil,
			salt,
		},
		{
			"invalid salt",
			pkey,
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewCredentials(tc.pkey, tc.salt)
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(errors.CausedBy(err, InvalidCredentialsCause)).To(gm.BeTrue())
		})
	}
}
