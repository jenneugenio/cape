package database

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestNoBackendSpecified(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("No backend specified", func(t *testing.T) {
		// I haven't specified a backend
		_, err := New("")
		gm.Expect(err).ToNot(gm.BeNil())
	})
}

func TestInvalidBackendSpecified(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("Invalid backend specified", func(t *testing.T) {
		_, err := New("fakedb://fake.db")
		_, ok := err.(*UnsupportedBackendError)

		gm.Expect(ok).To(gm.Equal(true))
	})
}

func TestValidBackend(t *testing.T) {
	gm.RegisterTestingT(t)
	t.Run("Valid backend specified", func(t *testing.T) {
		_, err := New("postgres://fake.db")
		gm.Expect(err).To(gm.BeNil())
	})
}
