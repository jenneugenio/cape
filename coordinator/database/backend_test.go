package database

import (
	"net/url"
	"testing"

	errors "github.com/capeprivacy/cape/partyerrors"
	gm "github.com/onsi/gomega"
)

func TestBackend(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("Invalid backend specified", func(t *testing.T) {
		u, err := url.Parse("fakedb://fake.db")
		gm.Expect(err).To(gm.BeNil())

		_, err = New(u, "test")
		gm.Expect(errors.FromCause(err, NotImplementedCause)).To(gm.BeTrue())
	})

	t.Run("Valid backend specified", func(t *testing.T) {
		u, err := url.Parse("postgres://fake.db")
		gm.Expect(err).To(gm.BeNil())

		_, err = New(u, "test")
		gm.Expect(err).To(gm.BeNil())
	})
}