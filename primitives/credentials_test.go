package primitives

import (
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestCredentials(t *testing.T) {
	gm.RegisterTestingT(t)

	creds, err := GenerateCredentials()
	gm.Expect(err).To(gm.BeNil())

	tests := []struct {
		name   string
		secret *base64.Value
		salt   *base64.Value
		alg    CredentialsAlgType

		cause *errors.Cause
	}{
		{
			name:   "valid credentials",
			secret: creds.Secret,
			salt:   creds.Salt,
			alg:    creds.Alg,
		},
		{
			name:   "invalid private key",
			secret: nil,
			salt:   creds.Salt,
			alg:    creds.Alg,
			cause:  &InvalidCredentialsCause,
		},
		{
			name:   "invalid salt",
			secret: creds.Secret,
			salt:   nil,
			alg:    creds.Alg,
			cause:  &InvalidCredentialsCause,
		},
		{
			name:   "invalid alg",
			secret: creds.Secret,
			salt:   creds.Salt,
			cause:  &InvalidCredentialsCause,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := NewCredentials(tc.secret, tc.salt, tc.alg)
			if tc.cause != nil {
				gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
				return
			}

			gm.Expect(out.Secret).To(gm.Equal(tc.secret))
			gm.Expect(out.Salt).To(gm.Equal(tc.salt))
			gm.Expect(out.Alg).To(gm.Equal(tc.alg))
		})
	}
}
