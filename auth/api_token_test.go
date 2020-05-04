package auth

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/primitives"
)

func TestAPIToken(t *testing.T) {
	gm.RegisterTestingT(t)

	host := "https://my.coordinator.com"
	u, err := primitives.NewURL(host)
	gm.Expect(err).To(gm.BeNil())

	userID, err := database.GenerateID(primitives.UserType)
	gm.Expect(err).To(gm.BeNil())

	secret, err := RandomSecret()
	gm.Expect(err).To(gm.BeNil())

	creds, err := NewCredentials(secret, nil)
	gm.Expect(err).To(gm.BeNil())

	pCreds, err := creds.Package()
	gm.Expect(err).To(gm.BeNil())

	tc, err := primitives.NewToken(userID, pCreds)
	gm.Expect(err).To(gm.BeNil())

	t.Run("new api token", func(t *testing.T) {
		token, err := NewAPIToken(secret, tc.ID, u)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.TokenID).To(gm.Equal(tc.ID))
		gm.Expect(token.URL).To(gm.Equal(u))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))
	})

	t.Run("marhsal unmarhal token", func(t *testing.T) {
		token, err := NewAPIToken(secret, tc.ID, u)
		gm.Expect(err).To(gm.BeNil())

		tokenStr, err := token.Marshal()
		gm.Expect(err).To(gm.BeNil())

		otherToken := &APIToken{}

		err = otherToken.Unmarshal(tokenStr)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(otherToken.TokenID).To(gm.Equal(tc.ID))
		gm.Expect(otherToken.URL).To(gm.Equal(u))
		gm.Expect(otherToken.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(otherToken.Secret)).To(gm.Equal(secretBytes))
		gm.Expect(otherToken.Secret).To(gm.Equal(token.Secret))

		gm.Expect(otherToken.Validate()).To(gm.BeNil())
	})

	t.Run("test unmarshal on raw string", func(t *testing.T) {
		writeToken, err := NewAPIToken(secret, tc.ID, u)
		gm.Expect(err).To(gm.BeNil())

		tokenStr, err := writeToken.Marshal()
		gm.Expect(err).To(gm.BeNil())

		token := &APIToken{}
		err = token.Unmarshal(tokenStr)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.TokenID).To(gm.Equal(tc.ID))
		gm.Expect(token.URL.String()).To(gm.Equal(host))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))

		gm.Expect(token.Validate()).To(gm.BeNil())
	})
}
