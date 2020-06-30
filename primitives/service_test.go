package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	gm.RegisterTestingT(t)

	email, err := NewEmail("service@cape.com")
	gm.Expect(err).To(gm.BeNil())

	t.Run("valid user service type", func(t *testing.T) {
		_, err := NewService(email, UserServiceType)
		gm.Expect(err).To(gm.BeNil())
	})
}
