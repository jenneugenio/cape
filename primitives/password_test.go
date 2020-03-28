package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"

	errors "github.com/dropoutlabs/cape/partyerrors"
)

func TestPassword(t *testing.T) {
	gm.RegisterTestingT(t)

	tests := map[string]struct {
		in    Password
		cause errors.Cause
	}{
		"short password": {
			in:    Password("abcd"),
			cause: InvalidPasswordCause,
		},
		"long password": {
			in:    Password("sdfkjsfkljasflkjs;fj;aksjfkalsjfkljslkjfalkjsdfjjlkfjkljsadf"),
			cause: InvalidPasswordCause,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.in.Validate()
			gm.Expect(errors.FromCause(err, test.cause)).To(gm.BeTrue())
		})
	}

	t.Run("passes for valid password", func(t *testing.T) {
		p := Password("abcdefgh")
		gm.Expect(p.Validate()).To(gm.BeNil())
	})

	t.Run("returns valid bytes", func(t *testing.T) {
		p := Password("helloeveryone")
		out := []byte("helloeveryone")

		gm.Expect(p.Bytes()).To(gm.Equal(out))
	})
}
