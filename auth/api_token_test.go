package auth

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/dropoutlabs/cape/primitives"
)

func TestAPIToken(t *testing.T) {
	gm.RegisterTestingT(t)

	email, err := primitives.NewEmail("email@email.com")
	gm.Expect(err).To(gm.BeNil())

	host := "https://my.controller.com"
	u, err := primitives.NewURL(host)
	gm.Expect(err).To(gm.BeNil())

	t.Run("new api token", func(t *testing.T) {
		token, err := NewAPIToken(email, u)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.Email).To(gm.Equal(email))
		gm.Expect(token.URL).To(gm.Equal(u))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))
	})

	t.Run("marhsal unmarhal token", func(t *testing.T) {
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

		gm.Expect(otherToken.Validate()).To(gm.BeNil())
	})

	t.Run("test unmarshal on raw string", func(t *testing.T) {
		tokenStr := "email@email.com,AQCiZ3kSIRgctnHV66K-SnxodHRwczovL215LmNvbnRyb2xsZXIuY29t"

		token := &APIToken{}
		err := token.Unmarshal(tokenStr)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.Email).To(gm.Equal(email))
		gm.Expect(token.URL.String()).To(gm.Equal(host))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))

		gm.Expect(token.Validate()).To(gm.BeNil())
	})
}
