package auth

import (
	"github.com/dropoutlabs/cape/primitives"
	"net/url"
	"testing"

	gm "github.com/onsi/gomega"
)

func TestAPIToken(t *testing.T) {
	gm.RegisterTestingT(t)

	email, err := primitives.NewEmail("email@email.com")
	gm.Expect(err).To(gm.BeNil())

	t.Run("new api token", func(t *testing.T) {
		host := "host.controller.com"

		u, err := url.Parse(host)
		gm.Expect(err).To(gm.BeNil())

		token, err := NewAPIToken(email, u)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.Email).To(gm.Equal(email))
		gm.Expect(token.URL).To(gm.Equal(u))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))
	})

	t.Run("marhsal unmarhal token", func(t *testing.T) {
		host := "host.controller.com"

		u, err := url.Parse(host)
		gm.Expect(err).To(gm.BeNil())

		token, err := NewAPIToken(email, u)
		gm.Expect(err).To(gm.BeNil())

		tokenStr, err := token.Marshal()
		gm.Expect(err).To(gm.BeNil())

		otherToken := &APIToken{}

		err = otherToken.Unmarshal(tokenStr)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(otherToken.Email).To(gm.Equal(email))
		gm.Expect(otherToken.URL).To(gm.Equal(u))
		gm.Expect(otherToken.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(otherToken.Secret)).To(gm.Equal(secretBytes))
		gm.Expect(otherToken.Secret).To(gm.Equal(token.Secret))
	})

	t.Run("test unmarshal on raw string", func(t *testing.T) {
		host := "host.controller.com"

		tokenStr := "email@email.com,AYqMLOkUUbK58Qr66G1a5v1ob3N0LmNvbnRyb2xsZXIuY29t"

		token := &APIToken{}
		err := token.Unmarshal(tokenStr)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.Email).To(gm.Equal(email))
		gm.Expect(token.URL.String()).To(gm.Equal(host))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))
	})
}
