package auth

import (
	"testing"

	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/models"
)

func TestSecret(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("can create secret from password", func(t *testing.T) {
		password := models.GeneratePassword()

		_, err := FromPassword(password)
		gm.Expect(err).To(gm.BeNil())
	})

	t.Run("can turn secret into password", func(t *testing.T) {
		password := models.GeneratePassword()

		secret, err := FromPassword(password)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(password).To(gm.Equal(secret.Password()))
	})
}

func TestAPIToken(t *testing.T) {
	gm.RegisterTestingT(t)

	userID := models.NewID()

	password := models.GeneratePassword()

	secret, err := FromPassword(password)
	gm.Expect(err).To(gm.BeNil())

	creds, err := DefaultSHA256Producer.Generate(password)
	gm.Expect(err).To(gm.BeNil())

	tc := models.NewToken(userID, &models.Credentials{
		Secret: creds.Secret,
		Salt:   creds.Salt,
		Alg:    creds.Alg,
	})
	gm.Expect(err).To(gm.BeNil())

	t.Run("new api token", func(t *testing.T) {
		token, err := NewAPIToken(secret, tc.ID)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.TokenID).To(gm.Equal(tc.ID))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))

		// token.Secret is the base64 value of the secret
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))
	})

	t.Run("marhsal unmarhal token", func(t *testing.T) {
		token, err := NewAPIToken(secret, tc.ID)
		gm.Expect(err).To(gm.BeNil())

		tokenStr, err := token.Marshal()
		gm.Expect(err).To(gm.BeNil())

		otherToken := &APIToken{}

		err = otherToken.Unmarshal(tokenStr)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(otherToken.TokenID).To(gm.Equal(tc.ID))
		gm.Expect(otherToken.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(otherToken.Secret)).To(gm.Equal(secretBytes))
		gm.Expect(otherToken.Secret).To(gm.Equal(token.Secret))

		gm.Expect(otherToken.Validate()).To(gm.BeNil())
	})

	t.Run("test unmarshal on raw string", func(t *testing.T) {
		writeToken, err := NewAPIToken(secret, tc.ID)
		gm.Expect(err).To(gm.BeNil())

		tokenStr, err := writeToken.Marshal()
		gm.Expect(err).To(gm.BeNil())

		token := &APIToken{}
		err = token.Unmarshal(tokenStr)
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(token.TokenID).To(gm.Equal(tc.ID))
		gm.Expect(token.Version).To(gm.Equal(uint8(tokenVersion)))
		gm.Expect(len(token.Secret)).To(gm.Equal(secretBytes))

		gm.Expect(token.Validate()).To(gm.BeNil())
	})
}
