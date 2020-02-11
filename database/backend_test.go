package database

import (
	"testing"

	errors "github.com/dropoutlabs/privacyai/partyerrors"
	gm "github.com/onsi/gomega"
)

func TestBackend(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("No backend specified", func(t *testing.T) {
		// I haven't specified a backend
		_, err := New("http://192.168.0.%31")
		gm.Expect(errors.FromCause(err, InvalidDBURLCause)).To(gm.BeTrue())
	})

	t.Run("Invalid backend specified", func(t *testing.T) {
		_, err := New("fakedb://fake.db")
		gm.Expect(errors.FromCause(err, NotImplementedDBCause)).To(gm.BeTrue())
	})

	t.Run("Valid backend specified", func(t *testing.T) {
		_, err := New("postgres://fake.db")
		gm.Expect(err).To(gm.BeNil())
	})
}
