package models

import (
	"crypto/rand"
	"testing"

	"github.com/manifoldco/go-base64"
	gm "github.com/onsi/gomega"
)

func TestPassword(t *testing.T) {
	gm.RegisterTestingT(t)

	pw := make([]byte, 128)
	_, err := rand.Read(pw)
	gm.Expect(err).To(gm.BeNil())

	tests := map[string]struct {
		in    Password
		cause string
	}{
		"short password": {
			in:    Password("abcd"),
			cause: "invalid_password: Passwords must be at least 8 characters long",
		},
		"long password": {
			in:    Password(base64.New(pw).String()),
			cause: "invalid_password: Passwords cannot be more than 128 characters long",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.in.Validate()
			gm.Expect(err).ToNot(gm.BeNil())
			gm.Expect(err.Error()).To(gm.Equal(test.cause))
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
