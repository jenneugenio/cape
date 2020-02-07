package database

import (
	"os"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestNoBackendSpecified(t *testing.T) {
	gm.RegisterTestingT(t)

	// I haven't specified a backend
	_, err := NewBackend()
	_, ok := err.(*UnspecifiedBackendError)

	gm.Expect(ok).To(gm.Equal(true))
}

func TestInvalidBackendSpecified(t *testing.T) {
	gm.RegisterTestingT(t)

	os.Setenv("DB_BACKEND", "fakedb")
	_, err := NewBackend()
	_, ok := err.(*UnsupportedBackendError)

	gm.Expect(ok).To(gm.Equal(true))
}

func TestValidBackend(t *testing.T) {
	gm.RegisterTestingT(t)

	os.Setenv("DB_BACKEND", "postgres")
	_, err := NewBackend()

	gm.Expect(err).To(gm.BeNil())
}
