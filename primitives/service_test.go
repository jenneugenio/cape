package primitives

import (
	"testing"

	gm "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	gm.RegisterTestingT(t)

	email, err := NewEmail("service@cape.com")
	gm.Expect(err).To(gm.BeNil())

	creds := &Credentials{}
	gm.Expect(err).To(gm.BeNil())

	endpoint, err := NewURL("https://service.cape.com")
	gm.Expect(err).To(gm.BeNil())

	t.Run("valid user service type", func(t *testing.T) {
		_, err := NewService(email, UserServiceType, nil, creds)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("valid data connector service type", func(t *testing.T) {
		_, err := NewService(email, DataConnectorServiceType, endpoint, creds)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("invalid user service type", func(t *testing.T) {
		_, err := NewService(email, UserServiceType, endpoint, creds)
		gm.Expect(err).ToNot(gm.BeNil())
	})

	t.Run("invalid data-connector service type", func(t *testing.T) {
		_, err := NewService(email, DataConnectorServiceType, nil, creds)
		gm.Expect(err).ToNot(gm.BeNil())
	})
}
