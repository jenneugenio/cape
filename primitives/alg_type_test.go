package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/capeprivacy/cape/partyerrors"
)

func TestAlgType(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := []struct {
		name  string
		alg   string
		cause *errors.Cause
	}{
		{
			name: "sha256 parses",
			alg:  "sha256",
		},
		{
			name: "argon2id parses",
			alg:  "argon2id",
		},
		{
			name: "eddsa parses",
			alg:  "eddsa",
		},
		{
			name:  "unknown algtype",
			alg:   "hello",
			cause: &InvalidAlgType,
		},
		{
			name:  "valid alg type - caps",
			alg:   "ARGON2ID",
			cause: &InvalidAlgType,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			alg := CredentialsAlgType(tc.alg)
			err := alg.Validate()
			if tc.cause != nil {
				gm.Expect(errors.FromCause(err, *tc.cause)).To(gm.BeTrue())
				return
			}

			gm.Expect(err).To(gm.BeNil())
		})
	}

	t.Run("Stringer", func(t *testing.T) {
		gm.Expect(EDDSA.String()).To(gm.Equal("eddsa"))
	})
}
